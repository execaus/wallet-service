package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToPgUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

func PgUUIDToUUID(p pgtype.UUID) (uuid.UUID, error) {
	if !p.Valid {
		return uuid.UUID{}, fmt.Errorf("pgtype.UUID is null")
	}
	return uuid.FromBytes(p.Bytes[:])
}
