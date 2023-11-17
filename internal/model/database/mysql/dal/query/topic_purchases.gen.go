// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/lixvyang/betxin.one/internal/model/database/mysql/dal/sqlmodel"
)

func newTopicPurchase(db *gorm.DB, opts ...gen.DOOption) topicPurchase {
	_topicPurchase := topicPurchase{}

	_topicPurchase.topicPurchaseDo.UseDB(db, opts...)
	_topicPurchase.topicPurchaseDo.UseModel(&sqlmodel.TopicPurchase{})

	tableName := _topicPurchase.topicPurchaseDo.TableName()
	_topicPurchase.ALL = field.NewAsterisk(tableName)
	_topicPurchase.ID = field.NewInt64(tableName, "id")
	_topicPurchase.TraceID = field.NewString(tableName, "trace_id")
	_topicPurchase.UID = field.NewString(tableName, "uid")
	_topicPurchase.Tid = field.NewInt64(tableName, "tid")
	_topicPurchase.YesPrice = field.NewString(tableName, "yes_price")
	_topicPurchase.NoPrice = field.NewString(tableName, "no_price")
	_topicPurchase.CreatedAt = field.NewInt64(tableName, "created_at")
	_topicPurchase.UpdatedAt = field.NewInt64(tableName, "updated_at")
	_topicPurchase.DeletedAt = field.NewInt64(tableName, "deleted_at")

	_topicPurchase.fillFieldMap()

	return _topicPurchase
}

type topicPurchase struct {
	topicPurchaseDo topicPurchaseDo

	ALL       field.Asterisk
	ID        field.Int64
	TraceID   field.String // 话题购买的trace_id
	UID       field.String
	Tid       field.Int64
	YesPrice  field.String // 支持金额
	NoPrice   field.String // 反对金额
	CreatedAt field.Int64
	UpdatedAt field.Int64
	DeletedAt field.Int64

	fieldMap map[string]field.Expr
}

func (t topicPurchase) Table(newTableName string) *topicPurchase {
	t.topicPurchaseDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t topicPurchase) As(alias string) *topicPurchase {
	t.topicPurchaseDo.DO = *(t.topicPurchaseDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *topicPurchase) updateTableName(table string) *topicPurchase {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.TraceID = field.NewString(table, "trace_id")
	t.UID = field.NewString(table, "uid")
	t.Tid = field.NewInt64(table, "tid")
	t.YesPrice = field.NewString(table, "yes_price")
	t.NoPrice = field.NewString(table, "no_price")
	t.CreatedAt = field.NewInt64(table, "created_at")
	t.UpdatedAt = field.NewInt64(table, "updated_at")
	t.DeletedAt = field.NewInt64(table, "deleted_at")

	t.fillFieldMap()

	return t
}

func (t *topicPurchase) WithContext(ctx context.Context) ITopicPurchaseDo {
	return t.topicPurchaseDo.WithContext(ctx)
}

func (t topicPurchase) TableName() string { return t.topicPurchaseDo.TableName() }

func (t topicPurchase) Alias() string { return t.topicPurchaseDo.Alias() }

func (t topicPurchase) Columns(cols ...field.Expr) gen.Columns {
	return t.topicPurchaseDo.Columns(cols...)
}

func (t *topicPurchase) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *topicPurchase) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 9)
	t.fieldMap["id"] = t.ID
	t.fieldMap["trace_id"] = t.TraceID
	t.fieldMap["uid"] = t.UID
	t.fieldMap["tid"] = t.Tid
	t.fieldMap["yes_price"] = t.YesPrice
	t.fieldMap["no_price"] = t.NoPrice
	t.fieldMap["created_at"] = t.CreatedAt
	t.fieldMap["updated_at"] = t.UpdatedAt
	t.fieldMap["deleted_at"] = t.DeletedAt
}

func (t topicPurchase) clone(db *gorm.DB) topicPurchase {
	t.topicPurchaseDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t topicPurchase) replaceDB(db *gorm.DB) topicPurchase {
	t.topicPurchaseDo.ReplaceDB(db)
	return t
}

type topicPurchaseDo struct{ gen.DO }

type ITopicPurchaseDo interface {
	gen.SubQuery
	Debug() ITopicPurchaseDo
	WithContext(ctx context.Context) ITopicPurchaseDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITopicPurchaseDo
	WriteDB() ITopicPurchaseDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITopicPurchaseDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITopicPurchaseDo
	Not(conds ...gen.Condition) ITopicPurchaseDo
	Or(conds ...gen.Condition) ITopicPurchaseDo
	Select(conds ...field.Expr) ITopicPurchaseDo
	Where(conds ...gen.Condition) ITopicPurchaseDo
	Order(conds ...field.Expr) ITopicPurchaseDo
	Distinct(cols ...field.Expr) ITopicPurchaseDo
	Omit(cols ...field.Expr) ITopicPurchaseDo
	Join(table schema.Tabler, on ...field.Expr) ITopicPurchaseDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITopicPurchaseDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITopicPurchaseDo
	Group(cols ...field.Expr) ITopicPurchaseDo
	Having(conds ...gen.Condition) ITopicPurchaseDo
	Limit(limit int) ITopicPurchaseDo
	Offset(offset int) ITopicPurchaseDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITopicPurchaseDo
	Unscoped() ITopicPurchaseDo
	Create(values ...*sqlmodel.TopicPurchase) error
	CreateInBatches(values []*sqlmodel.TopicPurchase, batchSize int) error
	Save(values ...*sqlmodel.TopicPurchase) error
	First() (*sqlmodel.TopicPurchase, error)
	Take() (*sqlmodel.TopicPurchase, error)
	Last() (*sqlmodel.TopicPurchase, error)
	Find() ([]*sqlmodel.TopicPurchase, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*sqlmodel.TopicPurchase, err error)
	FindInBatches(result *[]*sqlmodel.TopicPurchase, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*sqlmodel.TopicPurchase) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITopicPurchaseDo
	Assign(attrs ...field.AssignExpr) ITopicPurchaseDo
	Joins(fields ...field.RelationField) ITopicPurchaseDo
	Preload(fields ...field.RelationField) ITopicPurchaseDo
	FirstOrInit() (*sqlmodel.TopicPurchase, error)
	FirstOrCreate() (*sqlmodel.TopicPurchase, error)
	FindByPage(offset int, limit int) (result []*sqlmodel.TopicPurchase, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITopicPurchaseDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t topicPurchaseDo) Debug() ITopicPurchaseDo {
	return t.withDO(t.DO.Debug())
}

func (t topicPurchaseDo) WithContext(ctx context.Context) ITopicPurchaseDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t topicPurchaseDo) ReadDB() ITopicPurchaseDo {
	return t.Clauses(dbresolver.Read)
}

func (t topicPurchaseDo) WriteDB() ITopicPurchaseDo {
	return t.Clauses(dbresolver.Write)
}

func (t topicPurchaseDo) Session(config *gorm.Session) ITopicPurchaseDo {
	return t.withDO(t.DO.Session(config))
}

func (t topicPurchaseDo) Clauses(conds ...clause.Expression) ITopicPurchaseDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t topicPurchaseDo) Returning(value interface{}, columns ...string) ITopicPurchaseDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t topicPurchaseDo) Not(conds ...gen.Condition) ITopicPurchaseDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t topicPurchaseDo) Or(conds ...gen.Condition) ITopicPurchaseDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t topicPurchaseDo) Select(conds ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t topicPurchaseDo) Where(conds ...gen.Condition) ITopicPurchaseDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t topicPurchaseDo) Order(conds ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t topicPurchaseDo) Distinct(cols ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t topicPurchaseDo) Omit(cols ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t topicPurchaseDo) Join(table schema.Tabler, on ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t topicPurchaseDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t topicPurchaseDo) RightJoin(table schema.Tabler, on ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t topicPurchaseDo) Group(cols ...field.Expr) ITopicPurchaseDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t topicPurchaseDo) Having(conds ...gen.Condition) ITopicPurchaseDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t topicPurchaseDo) Limit(limit int) ITopicPurchaseDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t topicPurchaseDo) Offset(offset int) ITopicPurchaseDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t topicPurchaseDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITopicPurchaseDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t topicPurchaseDo) Unscoped() ITopicPurchaseDo {
	return t.withDO(t.DO.Unscoped())
}

func (t topicPurchaseDo) Create(values ...*sqlmodel.TopicPurchase) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t topicPurchaseDo) CreateInBatches(values []*sqlmodel.TopicPurchase, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t topicPurchaseDo) Save(values ...*sqlmodel.TopicPurchase) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t topicPurchaseDo) First() (*sqlmodel.TopicPurchase, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*sqlmodel.TopicPurchase), nil
	}
}

func (t topicPurchaseDo) Take() (*sqlmodel.TopicPurchase, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*sqlmodel.TopicPurchase), nil
	}
}

func (t topicPurchaseDo) Last() (*sqlmodel.TopicPurchase, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*sqlmodel.TopicPurchase), nil
	}
}

func (t topicPurchaseDo) Find() ([]*sqlmodel.TopicPurchase, error) {
	result, err := t.DO.Find()
	return result.([]*sqlmodel.TopicPurchase), err
}

func (t topicPurchaseDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*sqlmodel.TopicPurchase, err error) {
	buf := make([]*sqlmodel.TopicPurchase, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t topicPurchaseDo) FindInBatches(result *[]*sqlmodel.TopicPurchase, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t topicPurchaseDo) Attrs(attrs ...field.AssignExpr) ITopicPurchaseDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t topicPurchaseDo) Assign(attrs ...field.AssignExpr) ITopicPurchaseDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t topicPurchaseDo) Joins(fields ...field.RelationField) ITopicPurchaseDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t topicPurchaseDo) Preload(fields ...field.RelationField) ITopicPurchaseDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t topicPurchaseDo) FirstOrInit() (*sqlmodel.TopicPurchase, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*sqlmodel.TopicPurchase), nil
	}
}

func (t topicPurchaseDo) FirstOrCreate() (*sqlmodel.TopicPurchase, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*sqlmodel.TopicPurchase), nil
	}
}

func (t topicPurchaseDo) FindByPage(offset int, limit int) (result []*sqlmodel.TopicPurchase, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t topicPurchaseDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t topicPurchaseDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t topicPurchaseDo) Delete(models ...*sqlmodel.TopicPurchase) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *topicPurchaseDo) withDO(do gen.Dao) *topicPurchaseDo {
	t.DO = *do.(*gen.DO)
	return t
}