package repository

import (
	"context"
	"database/sql"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/repository/mysql"

	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (d *DB) Create(ctx context.Context,category entity.Category)(entity.Category,error){
	const op=richerror.Op("mysqlcategory.Create")

	logger.Info("Creating category in database","user_id",category.UserID,"name",category.Name)


	query:=`INSERT INTO categories (user_id,name,type,color) VALUES(?, ?, ?, ?)`

	res,err:=d.conn.Conn().ExecContext(ctx,query,category.UserID,category.Name,category.Type.String(),category.Color)

	if err !=nil{
		logger.Error("Failed to create category","user_id",category.UserID,"name",category.Name,"error",err.Error())

		return  entity.Category{},richerror.New(op).WithErr(err).WithMessage("failed to create category").WithKind(richerror.KindUnexpected)
	}
	
	id,err:=res.LastInsertId()
	if err!=nil{
		logger.Error("Failed to get inserted id","error",err.Error())
		return entity.Category{},richerror.New(op).WithErr(err).WithMessage("failed to get inserted id").WithKind(richerror.KindUnexpected)
	}
	category.ID=uint(id)
	logger.Info("category created successfully","category_id",category.ID)
	return category,nil
}


func (d *DB)GetByID(ctx context.Context,categoryID uint)(entity.Category,error){
	const op=richerror.Op("mysqlcategory.GetByID")

	query:=`SELECT id ,user_id, name,color,created_at FROM categories WHERE id = ?`

	row :=d.conn.Conn().QueryRowContext(ctx,query,categoryID)

	category,err:=scanCategory(row)
	if err !=nil{
		if err==sql.ErrNoRows {
			logger.Debug("Category not found","category_id",categoryID)
			return entity.Category{}, richerror.New(op).
				WithErr(err).
				WithMessage("category not found").
				WithKind(richerror.KindNotFound)
		}
		logger.Error("Failed to get category by ID","category_id",categoryID,"error",err.Error())

		return  entity.Category{},richerror.New(op).WithErr(err).WithMessage("failed to get category").WithKind(richerror.KindUnexpected)
	}
	return  category ,nil

}


func (d *DB) GetByUserIDAndType(ctx context.Context,userID uint,catType entity.CategoryType)([]entity.Category,error){
	const op = richerror.Op("mysqlcategory.GetByUserIDAndType")

	logger.Debug("Getting categories by user and type", "user_id",userID,"type",catType.String())

	query:=`SELECT id , user_id, name , type , color ,created_at FROM categories WHERE user_id = ? AND type = ?`

	rows,err :=d.conn.Conn().QueryContext(ctx,query,userID,catType.String())

	if err !=nil{
		logger.Error("Failed to get categories by user and type","user_id",userID,"type",catType.String(),"error",err.Error())
		return  nil,richerror.New(op).
						WithErr(err).
						WithMessage("failed to get categories").
						WithKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	var categories []entity.Category

	for rows.Next(){
		category,err:=scanCategory(rows)
		if err!=nil{
			logger.Error("Failed to scan category","user_id",userID,"error",err.Error())
			return nil,richerror.New(op).WithErr(err).WithMessage("failed to scan category").WithKind(richerror.KindUnexpected)
		}
		categories=append(categories, category)
	}

	if err=rows.Err();err!=nil{
		logger.Error("Error iterating categories",
			"user_id", userID,
			"error", err.Error(),
		)
		return nil, richerror.New(op).
			WithErr(err).
			WithMessage("failed to iterate categories").
			WithKind(richerror.KindUnexpected)
	}
	return  categories,nil

}


func (d *DB)GetByUserID(ctx context.Context , userID uint)([]entity.Category,error){
	const op=richerror.Op("mysqlcategory.GetByUserID")
	
	logger.Debug("Getting all categories for user","user_id",userID)

	query:=`SELECT  id ,user_id, name , type , color , created_at FROM categories WHERE user_id= ? `

	rows,err:=d.conn.Conn().QueryContext(ctx,query,userID)

	if err !=nil {
		logger.Error("Failed to get categories by user","user_id",userID,"error",err.Error())
		return nil,richerror.New(op).WithErr(err).WithMessage("failed to get categories").WithKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	var categories []entity.Category

	for rows.Next() {
		category, err:=scanCategory(rows)

		if err !=nil{
			logger.Error("Failed to scan category",
				"user_id", userID,
				"error", err.Error(),
			)
			return nil, richerror.New(op).
				WithErr(err).
				WithMessage("failed to scan category").
				WithKind(richerror.KindUnexpected)			
		}
		categories=append(categories, category)
	}
	if err =rows.Err(); err!=nil{
		logger.Error("Error iterating categories","user_id",userID,"error",err.Error())

		return  nil,richerror.New(op).WithErr(err).WithMessage("failed to iterate categories").WithKind(richerror.KindUnexpected)
	}

	return  categories,nil

}

func (d *DB) Update(ctx context.Context, category entity.Category) (entity.Category, error) {

	const op=richerror.Op("mysqlcategory.Update")

	logger.Info("Updating category in database","category_id",category.ID,"user_id",category.UserID)

	query:=`UPDATE categories SET name = ? , color = ? WHERE  id= ?`

	_,err:=d.conn.Conn().ExecContext(ctx,query,category.Name,category.Color,category.ID)

	if err !=nil{
		logger.Error("Failed to update category","category_id",category.ID,"error",err.Error())
		return entity.Category{},richerror.New(op).WithErr(err).WithMessage("failed to update category").WithKind(richerror.KindUnexpected)
	}

	logger.Info("Category updated successfully", "category_id", category.ID)

	return category, nil

}

// func (d *DB) Delete(ctx context.Context, categoryID uint) error {
func (d *DB) Delete(ctx context.Context, categoryID uint) error {
	const op=richerror.Op("mysqlcategory.Delete")
	logger.Info("Deleting category from database","category_id",categoryID)

	query:=`DELETE FROM categories WHERE  id = ? `

	res,err:=d.conn.Conn().ExecContext(ctx,query,categoryID)

	if err !=nil {
		logger.Error("Failed to delete category","category_id",categoryID,"error",err.Error())

		return richerror.New(op).WithErr(err).WithMessage("failed to delete category").WithKind(richerror.KindUnexpected)
	}
	rowsAffected,err:=res.RowsAffected()

	if err != nil {
		logger.Error("Failed to get rows affected", "error", err.Error())
		return richerror.New(op).
			WithErr(err).
			WithMessage("failed to verify deletion").
			WithKind(richerror.KindUnexpected)
	}
	if rowsAffected ==0 {
		logger.Warn("No Category found to delete","category_id",categoryID)

		return richerror.New(op).WithMessage("category not found").WithKind(richerror.KindNotFound)
	}

	logger.Info("Category deleted successfully","category_id",categoryID)

	return nil
}


func (d *DB) CategoryHasTransactions(ctx context.Context, categoryID uint) (bool, error) {
		const op = richerror.Op("mysqlcategory.CategoryHasTransactions")
		logger.Debug("Checking if category has transactions","category_id",categoryID)
		query := `SELECT EXISTS(SELECT 1 FROM transactions WHERE category_id = ?)`

		var exists bool
		err:=d.conn.Conn().QueryRowContext(ctx,query,categoryID).Scan(&exists)

		if err!=nil{
			logger.Error("Failed to check category transactions",
				"category_id", categoryID,
				"error", err.Error(),
			)
			return false, richerror.New(op).
				WithErr(err).
				WithMessage("failed to check category usage").
				WithKind(richerror.KindUnexpected)
		}
		logger.Debug("Category transactions check complete","category_id",categoryID,"has_trasactions",exists)
		return exists,nil
}

func scanCategory(scanner mysql.Scanner)(entity.Category,error){
	var category entity.Category
	var typeStr string
	var createdAt []uint8

	err:=scanner.Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&typeStr,
		&category.Color,
		&createdAt,
	)

	if err !=nil{
		return entity.Category{},err
	}

	category.Type=entity.MapToCategoryType(typeStr)

	return category,nil
}
