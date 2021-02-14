package core

import (
	"context"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/go-utils/code"
	"github.com/linger1216/go-utils/convert"
	"github.com/linger1216/go-utils/log"
	"github.com/linger1216/jelly-doc/src/server/pb"
	"strings"
)

const (
	HeadCountKey = "X-Total-Count"
)

type BasisApiDBService struct {
	logger *log.Log
	db     *sqlx.DB
	ddl    *BasisApiDDL
}

func NewApiDBService(logger *log.Log, db *sqlx.DB) pb.BasisApiServer {
	server := &BasisApiDBService{logger, db, NewBasisApiDDL()}
	query := server.ddl.CreateTableDDL()
	if _, err := server.db.Exec(query); err != nil {
		panic(err)
	}
	indexes := server.ddl.IndexTableDDL()
	for _, v := range indexes {
		if _, err := server.db.Exec(v); err != nil {
			panic(err)
		}
	}
	return server
}

func (f *BasisApiDBService) Create(ctx context.Context, in *pb.CreateApiRequest) (*pb.CreateApiResponse, error) {
	_ = ctx
	query, args := f.ddl.Upsert(in.Apis...)
	f.logger.Debugf("%s\n", query)

	_, err := f.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for i := range in.Apis {
		ids = append(ids, in.Apis[i].Id)
	}
	return &pb.CreateApiResponse{Ids: ids}, nil
}

func (f *BasisApiDBService) Delete(ctx context.Context, in *pb.DeleteApiRequest) (*pb.EmptyResponse, error) {
	_ = ctx
	query, args := f.ddl.Delete(in.Ids...)
	f.logger.Debugf("%s\n", query)
	_, err := f.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (f *BasisApiDBService) Update(ctx context.Context, in *pb.UpdateApiRequest) (*pb.EmptyResponse, error) {
	_ = ctx
	query, args := f.ddl.Upsert(in.Apis...)
	f.logger.Debugf("%s\n", query)
	_, err := f.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (f *BasisApiDBService) List(ctx context.Context, in *pb.ListApiRequest) (*pb.ListApiResponse, error) {
	_ = ctx
	resp := &pb.ListApiResponse{}
	query, args := f.ddl.List(in)
	f.logger.Debugf("%s\n", query)
	if in.Header > 0 {
		count := int64(0)
		err := f.db.Get(&count, query, args...)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, code.ErrNotFound
		}
		resp.Headers = append(resp.Headers, &pb.KV{
			Key:   HeadCountKey,
			Value: convert.Int64ToString(count),
		})
	} else {
		ret, err := f.query(query, args...)
		if err != nil {
			return nil, err
		}
		resp.Apis = ret
		if len(resp.Apis) == 0 {
			return nil, code.ErrNotFound
		}
	}
	return resp, nil
}

func (f *BasisApiDBService) query(query string, args ...interface{}) ([]*pb.Api, error) {
	rows, err := f.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	ret := make([]*pb.Api, 0)
	for rows.Next() {
		line := make(map[string]interface{})
		err = rows.MapScan(line)
		if err != nil {
			return nil, err
		}
		if tc, err := transApi("", line); err == nil && tc != nil {
			ret = append(ret, tc)
		}
	}

	if len(ret) == 0 {
		return nil, code.ErrNotFound
	}
	return ret, nil
}

func transApi(prefix string, m map[string]interface{}) (*pb.Api, error) {
	ret := &pb.Api{}

	if v, ok := m[prefix+"id"]; ok {
		ret.Id = convert.ToString(v)
	}

	if v, ok := m[prefix+"name"]; ok {
		ret.Name = convert.ToString(v)
	}

	if v, ok := m[prefix+"description"]; ok {
		ret.Description = convert.ToString(v)
	}

	if v, ok := m[prefix+"member_id"]; ok {
		ret.MemberIds = strings.Split(convert.ToString(v), ",")
	}

	if v, ok := m[prefix+"method"]; ok {
		ret.Method = convert.ToString(v)
	}

	if v, ok := m[prefix+"url"]; ok {
		ret.Url = convert.ToString(v)
	}

	if v, ok := m[prefix+"timeout"]; ok {
		ret.Timeout = int32(convert.ToInt64(v))
	}

	if v, ok := m[prefix+"headers"]; ok {
		if buf := convert.ToString(v); len(buf) > 0 {
			var m map[string]string
			err := jsoniter.ConfigFastest.UnmarshalFromString(buf, &m)
			if err != nil {
				return nil, err
			}
			ret.Headers = m
		}
	}

	if v, ok := m[prefix+"url_params"]; ok {
		if buf := convert.ToString(v); len(buf) > 0 {
			var m map[string]string
			err := jsoniter.ConfigFastest.UnmarshalFromString(buf, &m)
			if err != nil {
				return nil, err
			}
			ret.UrlParams = m
		}
	}

	if v, ok := m[prefix+"path_params"]; ok {
		if buf := convert.ToString(v); len(buf) > 0 {
			var m map[string]string
			err := jsoniter.ConfigFastest.UnmarshalFromString(buf, &m)
			if err != nil {
				return nil, err
			}
			ret.PathParams = m
		}
	}

	// auth
	//if v, ok := m[prefix+"auth"]; ok {
	//	if buf := convert.ToString(v); len(buf) > 0 {
	//		pb.Api_BasicAuth{}
	//		err := json_pb.UnmarshalString(buf, ret.Auth)
	//		if err != nil {
	//			return nil, err
	//		}
	//		ret.PathParams = m
	//	}
	//}

	if v, ok := m[prefix+"body"]; ok {
		ret.Body = convert.ToString(v)
	}

	if v, ok := m[prefix+"directories"]; ok {
		ret.Directories = strings.Split(convert.ToString(v), ",")
	}

	if v, ok := m[prefix+"create_time"]; ok {
		ret.CreateTime = convert.ToInt64(v)
	}
	if v, ok := m[prefix+"update_time"]; ok {
		ret.UpdateTime = convert.ToInt64(v)
	}

	return ret, nil
}

func (f *BasisApiDBService) Get(ctx context.Context, in *pb.GetApiRequest) (*pb.GetApiResponse, error) {
	_ = ctx
	query, args := f.ddl.Get(in.Ids...)
	f.logger.Debugf("%s\n", query)
	ret, err := f.query(query, args...)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, code.ErrNotFound
	}
	return &pb.GetApiResponse{Apis: ret}, nil
}

func (f *BasisApiDBService) Close() error {
	return f.db.Close()
}
