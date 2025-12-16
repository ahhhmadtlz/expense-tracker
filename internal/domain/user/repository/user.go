package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	"github.com/ahhhmadtlz/expense-tracker/internal/repository/mysql"
)



const (
	OpRegister=richerror.Op("mysqluser.Register")
	OpGetUserByID          = richerror.Op("mysql.GetUserByID")
	OpGetUserByPhoneNumber = richerror.Op("mysql.GetUserByPhoneNumber")
	OpIsPhoneNumberUnique=richerror.Op("mysqluser.IsPhoneNumberUnique")
)


func (d *DB) IsPhoneNumberUnique(ctx context.Context,phonenumber string)(bool ,error){

	const op=OpIsPhoneNumberUnique

	var exists bool

	err:=d.conn.Conn().QueryRowContext(ctx,`SELECT EXISTS(SELECT 1 FROM  users WHERE phone_number = ?)`,phonenumber).Scan(&exists)

	if err!=nil{
		d.logger.Error("failed to check phone number uniqueness","phone_number",phonenumber,"error",err.Error())

		return false,richerror.New(op).WithMessage("cant scan query result").WithKind(richerror.KindUnexpected).WithErr(err)
	}
	  d.logger.Debug("Phone number uniqueness checked",
        "phone_number", phonenumber,
        "is_unique", !exists,
    )

	return !exists,nil
}


func (d *DB)RegisterUser(ctx context.Context,user entity.User)(entity.User,error){
	const op=OpRegister

	query:=`INSERT INTO users(name ,phone_number,password,role,created_at)VALUES(?,?,?,?,NOW())`

	res,err:=d.conn.Conn().ExecContext(ctx,query,user.Name,user.PhoneNumber,user.Password,user.Role.String())
	if err !=nil{
		if isDuplicateKeyError(err){
			d.logger.Warn("Duplicate phone number", "phone_number", user.PhoneNumber)

			return entity.User{},richerror.New(op).WithErr(err).WithMessage("phone number already exists").WithKind(richerror.KindInvalid)
		}

		d.logger.Error("Failed to register user",
			"phone_number", user.PhoneNumber,
			"error", err.Error(),
		)

		return entity.User{},richerror.New(op).WithErr(err).WithMessage("cant execute command").WithKind(richerror.KindUnexpected)
	}

	id,err:=res.LastInsertId()
	if err !=nil{

		d.logger.Error("Failed to get inserted id", "error", err.Error())

		return entity.User{},richerror.New(op).WithErr(err).WithMessage("failed to get inserted id").WithKind(richerror.KindUnexpected)
	}

	user.ID=uint(id)
	d.logger.Info("User registered successfully", "user_id", user.ID)

	return user,nil
}


func (d *DB)GetUserByPhoneNumber(ctx context.Context,phoneNumber string)(entity.User,error){
	const op=OpGetUserByPhoneNumber
	row :=d.conn.Conn().QueryRowContext(ctx, `select * from users where phone_number = ?`, phoneNumber)

	user, err := scanUser(row)

	if err !=nil{
		if err==sql.ErrNoRows{
			d.logger.Debug("user not found","phone_number",phoneNumber)
				return entity.User{}, richerror.New(op).
				WithErr(err).
				WithMessage("not found").
				WithKind(richerror.KindNotFound)
		}
		d.logger.Error("Failed to get user by phone number",
			"phone_number", phoneNumber,
			"error", err.Error(),
		)
		return entity.User{}, richerror.New(op).
			WithErr(err).
			WithMessage("not found ").
			WithKind(richerror.KindUnexpected)
	}
	return user,nil
}

func (d *DB) GetUserByID(ctx context.Context, userID uint) (entity.User, error) {
	const op = OpGetUserByID
	row := d.conn.Conn().QueryRowContext(ctx, `select * from users where id = ? `, userID)
	user, err := scanUser(row)
	
	if err != nil {
		if err == sql.ErrNoRows {
			d.logger.Debug("User not found", "user_id", userID)
			return entity.User{}, richerror.New(op).
				WithErr(err).
				WithMessage("not found ").
				WithKind(richerror.KindNotFound)
		}
		d.logger.Error("Failed to get user by ID",
			"user_id", userID,
			"error", err.Error(),
		)
		return entity.User{}, richerror.New(op).
			WithMessage("not found ").
			WithKind(richerror.KindUnexpected)
	}
	return user, nil
}

func scanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User
	var createdAt []uint8

	var roleStr string

	err := scanner.Scan(
      &user.ID,           // id
      &user.Name,         // name
      &user.PhoneNumber,  // phone_number
      &roleStr,           // role (ENUM stored as string)
      &user.Password,     // password
      &createdAt,         // created_at
    )
		
	if err != nil {
		return entity.User{}, err
	}
    
	user.Role=entity.MapToRoleEntity(roleStr)

	return user,err
}

func isDuplicateKeyError(err error) bool {
    return strings.Contains(err.Error(), "Duplicate entry") ||
           strings.Contains(err.Error(), "1062")
}