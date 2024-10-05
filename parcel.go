package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	client := sql.Named("client", p.Client)
	status := sql.Named("status", p.Status)
	address := sql.Named("address", p.Address)
	createdAt := sql.Named("created_at", p.CreatedAt)

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)", client, status, address, createdAt)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления записи в БД: %w", err)
	}

	id, _ := res.LastInsertId()
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, fmt.Errorf("ошибка обработки записи БД: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, fmt.Errorf("ошибка получения записи БД: %w", err)
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка обработки записи из БД: %w", err)
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по записям БД: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса: %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status", sql.Named("address", address), sql.Named("number", number), sql.Named("status", "registered"))
	if err != nil {
		return fmt.Errorf("ошибка внесения изменений в БД: %w", err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status", sql.Named("number", number), sql.Named("status", "registered"))
	if err != nil {
		return fmt.Errorf("ошибка удаления записи в БД: %w", err)
	}
	return nil
}
