package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	"github.com/ahhhmadtlz/expense-tracker/internal/repository/mysql"
)

func (d *DB) Create(ctx context.Context, tx entity.Transaction)(entity.Transaction,error){
	const op=richerror.Op("mysqltransaction.Create")

	logger.Info("Creating transaction in database","user_id",tx.UserID,"type",tx.Type.String())

	query:=`INSERT INTO transactions (user_id, category_id, type, amount,description, date) VALUES(?,?,?,?,?,?)`

	res,err:=d.conn.Conn().ExecContext(ctx,query,tx.UserID,tx.CategoryID,tx.Type.String(),tx.Amount,tx.Description,tx.Date)


	if err !=nil{
		logger.Error("Failed to create transaction","user_id",tx.UserID,"error",err.Error())
		return entity.Transaction{},richerror.New(op).WithErr(err).WithMessage("failed to create transaction").WithKind(richerror.KindUnexpected)
	}

	id ,err:=res.LastInsertId()

	if err!=nil{
		logger.Error("Failed to get inserted id","error",err.Error())
		return entity.Transaction{},richerror.New(op).WithErr(err).WithMessage("failed to get inserted id").WithKind(richerror.KindUnexpected)
	}
	tx.ID = uint(id)
	logger.Info("Transaction created successfully", "transaction_id", tx.ID)
	return tx, nil
}

func (d *DB)GetByID(ctx context.Context,transactionID uint)(entity.Transaction,error){
	const op=richerror.Op("mysqltransaction.GetByID")

	query:=`SELECT id ,user_id,category_id, type, amount,description,date,created_at,updated_at FROM transaction WHERE id = ?`

	row :=d.conn.Conn().QueryRowContext(ctx,query,transactionID)

	tx,err:=scanTransaction(row)

	if err!=nil{
		if err ==sql.ErrNoRows{
			logger.Debug("Transaction not found","transaction_id",transactionID)
			return  entity.Transaction{},richerror.New(op).WithErr(err).WithMessage("transaction not found").WithKind(richerror.KindNotFound)
		}
		logger.Error("Failed to get transaction by ID","tranaction_id",transactionID,"error",err.Error())

		return  entity.Transaction{},richerror.New(op).WithErr(err).WithMessage("failed to get transaction").WithKind(richerror.KindUnexpected)
	}

	return  tx,nil
}

func (d *DB) GetByUserID(ctx context.Context , userID uint, filters map[string]any)([]entity.Transaction,error){
	const op=richerror.Op("mysqltransaction.GetByUserID")
	logger.Debug("Getting transactions for user","user_id",userID)

	query:=`SELECT id ,user_id,category_id,type,amount,description,date,created_at,updated_at FROM transactions WHERE user_id = ?`

	args:=[]any{userID}

	if txType,ok:=filters["type"].(string);ok&&txType!=""{
		query +=`AND type= ?`
		args=append(args, txType)
	}

	if categoryID,ok:=filters["category_id"].(uint);ok&&categoryID>0{
		query +=`AND category_id= ?`
		args=append(args, categoryID)
	}
	if startDate,ok:=filters["start_date"].(time.Time);ok&&!startDate.IsZero(){
		query+=`AND date >=?`
		args = append(args, startDate)
	}

	if endDate,ok:=filters["end_date"].(time.Time);ok&&!endDate.IsZero(){
		query +=`AND date <= ?`
		args=append(args, endDate)
	}

	query+=`ORDER BY date DESC, created_at DESC`

	rows,err:=d.conn.Conn().QueryContext(ctx,query,args...)
	if err!=nil{
		logger.Error("Failed to get transactions","user_id",userID,"error",err.Error())

		return  nil,richerror.New(op).WithErr(err).WithMessage("Failed to get transactions").WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()

	var transactions []entity.Transaction

	for rows.Next(){
		tx,err:=scanTransaction(rows)
		if err!=nil{
			logger.Error("Failed to scan transaction","user_id",userID,"error",err.Error())
			return nil,richerror.New(op).WithErr(err).WithMessage("failed to scan transaction").WithKind(richerror.KindUnexpected)
		}
		transactions=append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error iterating transactions","user_id",userID,"error",err.Error())
		return  nil ,richerror.New(op).WithErr(err).WithMessage("failed to iterate transactions").WithKind(richerror.KindUnexpected)
	}

	return  transactions,nil

}

func (d *DB) Update(ctx context.Context, tx entity.Transaction) (entity.Transaction, error) {
	const op = richerror.Op("mysqltransaction.Update")

	logger.Info("Updating transaction in database", "transaction_id", tx.ID, "user_id", tx.UserID)

	query := `UPDATE transactions SET category_id = ?, type = ?, amount = ?, description = ?, date = ? WHERE id = ?`

	_, err := d.conn.Conn().ExecContext(ctx, query, tx.CategoryID, tx.Type.String(), tx.Amount, tx.Description, tx.Date, tx.ID)

	if err != nil {
		logger.Error("Failed to update transaction", "transaction_id", tx.ID, "error", err.Error())
		return entity.Transaction{}, richerror.New(op).WithErr(err).WithMessage("failed to update transaction").WithKind(richerror.KindUnexpected)
	}

	logger.Info("Transaction updated successfully", "transaction_id", tx.ID)
	return tx, nil
}

func (d *DB) Delete(ctx context.Context,transactionID uint) error {
	const op=richerror.Op("mysqltransaction.Delete")
	
	logger.Info("Deleting transaction from database","transaction_id",transactionID)

	query:=`DELETE FROM transactions WHERE id= ?`

	res,err:=d.conn.Conn().ExecContext(ctx,query,transactionID)

	if err!=nil{
		logger.Error("Failed to delete transaction","transaction_id",transactionID,"error",err.Error())
		return richerror.New(op).WithErr(err).WithMessage("failed to delete transaction").WithKind(richerror.KindUnexpected) 
	}
	
	rowsAffected,err:=res.RowsAffected()
	if err != nil {
		logger.Error("Failed to get rows affected", "error", err.Error())
		return richerror.New(op).WithErr(err).WithMessage("failed to verify deletion").WithKind(richerror.KindUnexpected)
	}

	if rowsAffected == 0 {
		logger.Warn("No transaction found to delete", "transaction_id", transactionID)
		return richerror.New(op).WithMessage("transaction not found").WithKind(richerror.KindNotFound)
	}

	logger.Info("Transaction deleted successfully", "transaction_id", transactionID)
	return nil
}

func scanTransaction(scanner mysql.Scanner)(entity.Transaction,error){
	var tx entity.Transaction
	var typeStr string
	var date,createdAt,updatedAt []uint8
	err:=scanner.Scan(
		&tx.ID,
		&tx.UserID,
		&tx.CategoryID,
		&typeStr,
		&tx.Amount,
		&tx.Description,
		&date,
		&createdAt,
		&updatedAt,
	)

	if err!=nil{
		return  entity.Transaction{},err
	}

	tx.Type=entity.MapToTransactionType(typeStr)

	if len(date)>0 {
		tx.Date,_=time.Parse("2006-01-02 15:04:05",string(date))
	}
	if len(createdAt)>0{
		tx.CreatedAt,_=time.Parse("2006-01-02 15:04:05",string(createdAt))
	}

	if len(updatedAt)>0{
		tx.UpdatedAt,_=time.Parse("2006-01-02 15:04:05",string(updatedAt))
	}

	return tx,nil

}