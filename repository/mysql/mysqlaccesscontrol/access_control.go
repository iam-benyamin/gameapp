package mysqlaccesscontrol

import (
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"gameapp/pkg/slice"
	"gameapp/repository/mysql"
	"strings"
	"time"
)

func (d *DB) GetUserPermissionTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error) {
	const op = "mysql.GetUserPermissionTitles"
	//user, err := d.GetUserByID(userID)
	//if err != nil {
	//	return nil, richerror.New(op).WithErr(err)
	//}

	roleACL := make([]entity.AccessControl, 0)

	rows, err := d.conn.Conn().Query(`SELECT * FROM access_controls WHERE actor_type = ? AND  actor_id = ?`,
		entity.RoleActorType, role)

	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
			WithKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := scanAccessControl(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
				WithKind(richerror.KindUnexpected)
		}

		roleACL = append(roleACL, acl)
	}

	if err = rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
			WithKind(richerror.KindUnexpected)
	}

	userACL := make([]entity.AccessControl, 0)

	userRows, err := d.conn.Conn().Query(`SELECT * FROM access_controls WHERE actor_type = ? AND  actor_id = ?`,
		entity.UserActorType, userID,
	)

	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
			WithKind(richerror.KindUnexpected)
	}

	defer userRows.Close()

	for userRows.Next() {
		acl, err := scanAccessControl(userRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
				WithKind(richerror.KindUnexpected)
		}

		userACL = append(userACL, acl)
	}

	if err = userRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
			WithKind(richerror.KindUnexpected)
	}

	// merge acl by permissions id
	permissionIDs := make([]uint, 0)
	for _, r := range roleACL {
		if !slice.DoseExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}
	}

	if len(permissionIDs) == 0 {
		return nil, nil
	}

	// select * from permissions where id in (1, 3, 7)
	args := make([]interface{}, len(permissionIDs))

	for i, id := range permissionIDs {
		args[i] = id
	}

	// this query works if we have one or more permissions id
	query := "select * from permissions where id in (?" + strings.Repeat(",?", len(permissionIDs)-1) + ")"

	pRows, err := d.conn.Conn().Query(query, args...)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
			WithKind(richerror.KindUnexpected)
	}
	defer pRows.Close()

	permissionTitles := make([]entity.PermissionTitle, 0)

	for pRows.Next() {
		permission, err := scanPermission(pRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
				WithKind(richerror.KindUnexpected)
		}

		permissionTitles = append(permissionTitles, permission.Title)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.ErrorSomeThingWentWrong).
			WithKind(richerror.KindUnexpected)
	}

	return permissionTitles, nil
}

func scanAccessControl(scanner mysql.Scanner) (entity.AccessControl, error) {
	var createdAt time.Time
	var acl entity.AccessControl

	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &createdAt)

	return acl, err
}
