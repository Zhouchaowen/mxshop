package handler

import (
	"context"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
	"sync"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	// 设置库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}

	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// TODO
func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存 本地事务，数据一致性
	tx := global.DB.Begin()
	// 有Clauses(clause.Locking{Strength: "UPDATE"})就可以不要mx.lock()
	//mx.Lock() // 性能问题，锁范围太大，分布式锁
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		// Clauses(clause.Locking{Strength: "UPDATE"})数据库锁 for update
		// 悲观锁实现
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback()
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}

		// 乐观锁实现
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减
		// TODO 并发可能会超卖，数据不一致。锁，分布式锁。
		inv.Stocks -= goodInfo.Num

		//tx.Save(&inv)

		// update inventory set stocks = stocks-1,version=version+1 where goods=goods and version=version
		// 这种写法有坑，零值 对于int类型来说 默认为0，不更新
		//if result := tx.Model(&model.Inventory{}).Where("goods = ? and version = ?",goodInfo.GoodsId,inv.Version).Updates(model.Inventory{Stocks: inv.Stocks,Version: inv.Version});result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//}else {
		//	break
		//}

		// 乐观锁实现
		// 零值 对于int类型来说 默认为0，不更新 通过Select来解决
		if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version = ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version}); result.RowsAffected == 0 {
			zap.S().Info("库存扣减失败")
		} else {
			break
		}
	}
	tx.Commit() // 手动提交修改
	//mx.Unlock()
	return &emptypb.Empty{}, nil
}

// 本地锁实现方式
var mx sync.Mutex

func (*InventoryServer) Sell1(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存 本地事务，数据一致性
	tx := global.DB.Begin()
	mx.Lock() // 性能问题，锁范围太大，分布式锁
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 手动提交修改
	mx.Unlock()
	return &emptypb.Empty{}, nil
}

// 数据库悲观锁实现方式
func (*InventoryServer) Sell2(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存 本地事务，数据一致性
	tx := global.DB.Begin()
	// 有Clauses(clause.Locking{Strength: "UPDATE"})就可以不要mx.lock()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		// Clauses(clause.Locking{Strength: "UPDATE"})数据库锁 for update
		// 悲观锁实现
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减
		// TODO 并发可能会超卖，数据不一致。锁，分布式锁。
		inv.Stocks -= goodInfo.Num

		tx.Save(&inv)
	}
	tx.Commit() // 手动提交修改
	return &emptypb.Empty{}, nil
}

// 数据库乐观锁实现方式
func (*InventoryServer) Sell3(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 扣减库存 本地事务，数据一致性
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		// 乐观锁实现
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "没有库存信息")
		}

		if inv.Stocks < goodInfo.Num {
			tx.Rollback() // 回滚操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		// 扣减
		inv.Stocks -= goodInfo.Num

		// update inventory set stocks = stocks-1,version=version+1 where goods=goods and version=version
		// 这种写法有坑，零值 对于int类型来说 默认为0，不更新
		//if result := tx.Model(&model.Inventory{}).Where("goods = ? and version = ?",goodInfo.GoodsId,inv.Version).Updates(model.Inventory{Stocks: inv.Stocks,Version: inv.Version});result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//}else {
		//	break
		//}

		// 乐观锁实现
		// 零值 对于int类型来说 默认为0，不更新 通过Select来解决
		if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").
			Where("goods = ? and version = ?", goodInfo.GoodsId, inv.Version).
			Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version}); result.RowsAffected == 0 {
			zap.S().Info("库存扣减失败")
		} else {
			break
		}
	}
	tx.Commit() // 手动提交修改
	return &emptypb.Empty{}, nil
}

/**
分布式锁要解决的问题？
	a.互斥
	b.死锁
	c.安全性，只能被持有该锁的用户删除，删除时会对比设置进去的value。
*/
// redis实现分布式锁形式
/**
原理：setnx

问题：
	1.设置过期时间
	2.过期时间到达，业务未执行完。
		解决：定期刷新 内部通过lua脚本实现，需要自己去启动协程完成续租刷新。
	3.redis在集群下同步数据延迟，导致分布式锁的问题。redlock
		在分布式环境下，异步拿锁，只要拿到n/2+1的锁，就加锁成功。
		拿到锁后计算过期时间


*/
func (*InventoryServer) Sell4(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.0.104:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	tx := global.DB.Begin()

	//这个时候应该先查询表，然后确定这个订单是否已经扣减过库存了，已经扣减过了就别扣减了
	//并发时候会有漏洞， 同一个时刻发送了重复了多次， 使用锁，分布式锁
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		// 设置锁定的值
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		// TODO 获取分布式锁
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)

		// TODO 释放分布式锁
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}

	}
	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

// TODO
func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// 规划库存，1.订单超时归还，2.订单创建失败，归还扣减库存，3.手动归还
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}

		// 扣减
		// TODO 并发可能会超卖，数据不一致。锁，分布式锁。
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 手动提交修改
	return &emptypb.Empty{}, nil
}
