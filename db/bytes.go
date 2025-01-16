package db

import "zombiezen.com/go/sqlite"

func StmtBytes(stmt *sqlite.Stmt, colName string) []byte {
	bl := stmt.GetLen(colName)
	if bl == 0 {
		return nil
	}

	buf := make([]byte, bl)
	if writtent := stmt.GetBytes(colName, buf); writtent != bl {
		return nil
	}

	return buf
}

func StmtBytesByCol(stmt *sqlite.Stmt, col int) []byte {
	bl := stmt.ColumnLen(col)
	if bl == 0 {
		return nil
	}

	buf := make([]byte, bl)
	if writtent := stmt.ColumnBytes(col, buf); writtent != bl {
		return nil
	}

	return buf
}
