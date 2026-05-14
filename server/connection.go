package server

import (
	"bufio"
	"fmt"
	"go/parser"
	"net"
	"strata_db/engine"
	"strings"
)



type Connection struct {
	conn net.Conn
	db  *engine.Database
}


func NewConnection(conn net.Conn, db *engine.Database) *Connection {
	return &Connection{
		conn: conn,
		db: db,
	}
}


func (c *Connection) Handle() error {

	reader := bufio.NewReader(c.conn)
	
	// Send welcome message
	c.writeResult("Stratadb Ready\n")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			c.writePrompt()
			continue
		}
		if line == "EXIT" || line == "QUIT" {
			c.writeResult("Goodbye\n")
			return nil
		}

		// Parse and execute
		result, err := c.executeQuery(line)
		if err != nil {
			c.writeResult(fmt.Sprintf("Error: %v\n"))
		}else {
			c.writeResult(result)
		}
		c.writePrompt()
	}
}


func (c *Connection) executeQuery(query string) (string, error) {
	p := parser.NewParser(query)
	stmt, err := p.Parse()
	if err != nil {
		return  "", err
	}

	switch stmt.Type {
    case parser.StatementCreate:
        return c.handleCreate(stmt)
    case parser.StatementInsert:
        return c.handleInsert(stmt)
    case parser.StatementSelect:
        return c.handleSelect(stmt)
    case parser.StatementUpdate:
        return c.handleUpdate(stmt)
    case parser.StatementDelete:
        return c.handleDelete(stmt)
    default:
        return "", fmt.Errorf("unsupported statement type")
    }
}



func (c *Connection) executeQuery(query string) (string, error ) {}



func (c *Connection) handleCreate(stmt *parser.Statement) (string, error) {
	return c.db.CreateTable(stmt.TableName, stmt.Columns)
}


func (c *Connection) handleInsert(stmt *parser.Statement) (string, error) {
	return c.db.Insert(stmt.TableName, stmt.Values)
}

func (c *Connection) handleSelect(stmt *parser.Statement) (string, error) {
	return c.db.Select(stmt.TableName, stmt.Columns, stmt.Condition)
}

func (c *Connection) handleUpdate(stmt *parser.Statement) (string, error) {
	return c.db.Update(stmt.TableName, stmt.SetValues, stmt.Condition)
}

func (c Connection) handleDelete(stmt *parser.Statement) (string, error) {
	return c.db.Delete(stmt.TableName, stmt.Condition)
}

func (c *Connection) writeResult(msg string) {
	c.conn.Write([]byte(msg))
}

func (c *Connection) writePrompt(){
	c.conn.Write([]byte("> "))
}


