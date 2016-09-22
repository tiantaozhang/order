package common

import (
	"database/sql"
	"github.com/Centny/gwf/dbutil"
)

func CheckDb(db *sql.DB) error {
	_, err := dbutil.DbQueryI(db, `SELECT COUNT(*) FROM ods_order`)
	if err != nil {
		err = dbutil.DbExecScript(db, orderNew)
		// _, err = db.Exec(funcs)
	}
	return err
}

// var order_sql string = `
// /*==============================================================*/
// /* DBMS name:      MySQL 5.0                                    */
// /* Created on:     2015-09-11 17:01:59                          */
// /*==============================================================*/

// drop table if exists ODS_APP;

// drop table if exists ODS_BALANCE;

// drop table if exists ODS_CART;

// drop table if exists ODS_DISCOUNT;

// drop table if exists ODS_INCOME;

// drop table if exists ODS_ORDER;

// drop table if exists ODS_ORDER_ITEM;

// drop table if exists ODS_RECORD;

// drop table if exists ODS_REFUND;

// drop table if exists ODS_REFUND_ITEM;

// drop table if exists ODS_SELLER;

// drop table if exists ODS_EVENT;

// /*==============================================================*/
// /* Table: ODS_APP                                               */
// /*==============================================================*/
// create table ODS_APP
// (
//    TID                  int not null auto_increment,
//    NAME                 varchar(255),
//    APP_ID               varchar(255),
//    APP_ICON             varchar(255),
//    DETAIL_URL           varchar(255) comment '详情回调',
//    REFUND_URL           varchar(255) comment '退款回调',
//    REDICT_URL           varchar(255) comment '回调',
//    TOKEN                varchar(255),
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_BALANCE                                           */
// /*==============================================================*/
// create table ODS_BALANCE
// (
//    TID                  int not null auto_increment,
//    BALANCE              float,
//    USER_ID              int,
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_CART                                              */
// /*==============================================================*/
// create table ODS_CART
// (
//    TID                  int not null auto_increment,
//    APP_ID               varchar(255) not null comment '订单来源系统ID号',
//    GOODS_ID             varchar(255) comment '商品ID',
//    GOODS_DETAIL         varchar(255) comment '商品详情信息（json数据）',
//    SELLER               int not null comment '卖家',
//    SELLER_NAME          varchar(255) comment '卖家名称',
//    BUYER                int not null comment '买家',
//    BUYER_NAME           varchar(255) comment '买家名称',
//    COUNT                int not null comment '数量',
//    PRICE                decimal(16,4) not null comment '单价',
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态(N:正常,D:删除状态,I:直接购买)',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_DISCOUNT                                          */
// /*==============================================================*/
// create table ODS_DISCOUNT
// (
//    TID                  int not null auto_increment,
//    APP_ID               varchar(255) comment '来源系统',
//    USER                 int comment '创建用户ID',
//    DISCOUNT_NUM         int comment '优惠券数量',
//    COURSE_OR_SELLER_ID  varchar(255) comment '课程或教师ID 可多个',
//    PROFIT_SYS_PERCENT   float(13,2) comment '系统百分比提成',
//    PROFIT_SELLER_PERCENT float(13,2) comment '销售百分比提成',
//    PROFIT_FROM_SYS      tinyint comment '1:销售百分比提成从系统百分比提成中扣除 0:不扣除',
//    PRICE_LIMIT          varchar(255) comment '价格限制 可多个',
//    DISCOUNT_PRICE       varchar(255) comment '减的价格或折扣 可多个',
//    DISCOUNT_TYPE        varchar(255) comment '优惠类型(REDUCE满减，PERCENT满多少折扣多少)',
//    TYPE                 varchar(255) comment '类型(ALL通用   COURSE课程   TEACHER老师)',
//    START_TIME           timestamp default '0000-00-00 00:00:00' comment '开始时间',
//    END_TIME             timestamp comment '结束时间',
//    NAME                 varchar(255) comment '优惠名称',
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_INCOME                                            */
// /*==============================================================*/
// create table ODS_INCOME
// (
//    TID                  int not null auto_increment,
//    SELLER_ID            int,
//    ORDER_NO             varchar(255),
//    ORDER_ID             int,
//    MONEY                float,
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_ORDER                                             */
// /*==============================================================*/
// create table ODS_ORDER
// (
//    TID                  int not null auto_increment,
//    ORD_NO               varchar(255) not null comment '订单号',
//    RES_APP              varchar(255) comment '订单来源系统',
//    SELLER               int comment '卖家',
//    BUYER                int not null comment '买家',
//    TOTAL_PRICE          decimal(16,4) not null comment '总价',
//    PRICE                decimal(16,4) not null comment '总价',
//    SELLER_NAME          varchar(255) comment '卖家名称',
//    PAY_WAY              varchar(255) comment '支付方式',
//    BUYER_NAME           varchar(255) comment '买家名称',
//    ORD_TYPE             varchar(255) not null comment '订单类型
//             SHOPPING,RECHARGE,REF',
//    CHANNEL              varchar(255) comment '渠道',
//    DISCOUNT             varchar(255) comment '优惠',
//    PAY_TIME             timestamp default '0000-00-00 00:00:00' comment '支付时间',
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_ORDER_ITEM                                        */
// /*==============================================================*/
// create table ODS_ORDER_ITEM
// (
//    TID                  int not null auto_increment,
//    ORDER_ID             int not null comment '订单id',
//    ORD_NO               varchar(255) not null comment '订单号',
//    GOODS_ID             varchar(255) comment '商品id',
//    GOODS_DETAIL         varchar(255),
//    PRICE                decimal(16,4) comment '商品单价',
//    TOTAL_PRICE          decimal(16,4) not null comment '总价',
//    DISCOUNT_PRICE       decimal(16,4) not null comment '优惠后总价',
//    BUYER                int not null comment '买家',
//    SELLER               int comment '卖家',
//    COUNT                int comment '购买数量',
//    COMMENT_STATUS       varchar(255) comment '评论状态',
//    REMARK               varchar(255) comment '备注',
//    TIME                 timestamp comment '下单时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_RECORD                                            */
// /*==============================================================*/
// create table ODS_RECORD
// (
//    TID                  int not null auto_increment,
//    USER_ID              int,
//    TARGET_ID            varchar(255),
//    NAME                 varchar(255),
//    MONEY                float,
//    TYPE                 varchar(255),
//    DESCRIPTION          varchar(255),
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_REFUND                                            */
// /*==============================================================*/
// create table ODS_REFUND
// (
//    TID                  int not null auto_increment,
//    ORDER_ID             int not null comment '订单号',
//    ORDER_NO             varchar(255) not null comment '订单号',
//    REFUND_NO            varchar(255) not null comment '订单号',
//    REASON               varchar(255) not null comment '订单号',
//    RES_APP              varchar(255) comment '订单来源系统',
//    SELLER               int comment '卖家',
//    BUYER                int not null comment '买家',
//    PRICE                decimal(16,4) not null comment '总价',
//    SELLER_NAME          varchar(255) comment '卖家名称',
//    BUYER_NAME           varchar(255) comment '买家名称',
//    PAY_WAY              varchar(255) comment '支付方式',
//    ORD_TYPE             varchar(255) not null comment '订单类型
//             SHOPPING,RECHARGE,REF',
//    REFUND_TIME          timestamp default '0000-00-00 00:00:00' comment '支付时间',
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_REFUND_ITEM                                       */
// /*==============================================================*/
// create table ODS_REFUND_ITEM
// (
//    TID                  int not null auto_increment,
//    REFUND_ID            int not null comment '订单id',
//    ORD_NO               varchar(255) not null comment '订单号',
//    RES_ID               int comment '商品id',
//    RES_TYPE             int comment '商品类型 1课程 2资源',
//    PRICE                decimal(16,4) comment '商品单价',
//    BUYER                int not null comment '买家',
//    SELLER               int comment '卖家',
//    COUNT                int comment '购买数量',
//    REMARK               varchar(255) comment '备注',
//    TIME                 datetime comment '下单时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_SELLER                                            */
// /*==============================================================*/
// create table ODS_SELLER
// (
//    TID                  int not null auto_increment,
//    DISCOUNT_ID          int comment '优惠券ID',
//    SELLER_ID            int comment '销售ID',
//    CODE                 varchar(255) comment '优惠代码',
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255),
//    primary key (TID)
// );

// /*==============================================================*/
// /* Table: ODS_EVENT                                             */
// /*==============================================================*/
// create table ODS_EVENT
// (
//    TID                  int not null auto_increment PRIMARY KEY,
//    EVENT_ID             int,
//    EVENT_NAME           varchar(255),
//    EVENT_PARAM          BLOB,
//    EVENT_ACTION         varchar(255),
//    EVENT_STATUS         varchar(255),
//    EVENT_FAIL_TIMES     INT,
//    ERROR_MSG            varchar(2000),
//    TYPE                 varchar(255),
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255),
//    ADD2                 varchar(255)
// );

// /*==============================================================*/
// /* Table: ODS_ORDER_SET                                         */
// /*==============================================================*/
// create table ODS_ORDER_SET
// (
//    TID                  int not null auto_increment PRIMARY KEY,
//    TRADE_NO             varchar(255),
//    ORDER_NOS            varchar(255),
//    PAY_WAY              varchar(255),
//    MSG                  BLOB,
//    TIMES                INT,
//    TIME                 timestamp comment '创建时间',
//    STATUS               varchar(255) not null comment '状态',
//    ADD1                 varchar(255)
// );
// `

var orderNew string = `
/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2015/12/1 10:15:59                           */
/*==============================================================*/


drop table if exists ODS_ORDER;

drop table if exists ODS_ORDER_ENV;

drop table if exists ODS_ORDER_ITEM;

drop table if exists ODS_ORDER_REFUND;

drop table if exists ODS_RECORD;

/*==============================================================*/
/* Table: ODS_ORDER                                             */
/*==============================================================*/
create table ODS_ORDER
(
   tid                  int not null AUTO_INCREMENT,
   ono                  varchar(255) comment '订单编号',
   buyer                int comment '买家id',
   seller               int comment '卖家id',
   total_price          float comment '总价',
   type                 varchar(255) comment '类型
            N  正常状态
            REFUNG  退款状态',
   status               varchar(255) comment '状态
            未付款NOT_PAY
            已付款PAID
            停用INVALID
            ',
   time                 timestamp,
   return_url varchar(255) comment '网页回调地址',
   expand varchar(255) comment '额外参数',
   wno    varchar(255) comment '微信订单号',
   add1                 varchar(255),
   add2                 varchar(255),
   primary key (tid)
);

alter table ODS_ORDER comment '订单表';

/*==============================================================*/
/* Table: ODS_ORDER_ENV                                         */
/*==============================================================*/
create table ODS_ORDER_ENV
(
   tid                  int not null AUTO_INCREMENT,
   akey                 varchar(255) comment '当type是PAID_CB/REFUND_CB，对应order_item的from字',
   aval                 varchar(255) comment '回调地址以format形式保存，如https://sss.com/cb?ono=%s',
   type                 varchar(255) comment '类型',
   status               varchar(255) comment '状态',
   time                 timestamp,
   add1                 varchar(255),
   add2                 varchar(255),
   primary key (tid)
);

alter table ODS_ORDER_ENV comment '订单环境';

/*==============================================================*/
/* Table: ODS_ORDER_ITEM                                            */
/*==============================================================*/
create table ODS_ORDER_ITEM
(
   tid                  int not null AUTO_INCREMENT,
   ono                  varchar(255) comment '订单编号',
   oid                  int comment '订单id',
   p_name               varchar(255) comment '物品名称',
   p_id                 int comment '物品id',
   p_img                varchar(255) comment '物品图片',
   p_type               varchar(255) comment '物品类型',
   p_count              int comment '物品数量',
   p_from               varchar(255) comment '来源系统',
   notified             int comment '已通知次数',
   price                float comment '价格',
   type                 varchar(255) comment '类型
            N  正常
            REFUND 退款',
   status               varchar(255) comment '状态
            订单类型为退款时：REFUNDING_B(买家提交资料）20/REFUNDING_S（卖家提交资料）30/REFUNDING_SERVICE（客服）40/REFUNDING_CANCEL(取消）50/REFUNDED 60/INVALID -1
                  订单类型为正常时：N（10） REFUNDED 60/INVALID -1',
   time                 timestamp,
   add1                 varchar(255),
   add2                 varchar(255),
   primary key (tid)
);

alter table ODS_ORDER_ITEM comment '订单内容';

/*==============================================================*/
/* Table: ODS_ORDER_REFUND                                      */
/*==============================================================*/
create table ODS_ORDER_REFUND
(
   tid                  int not null AUTO_INCREMENT,
   ono                  varchar(255) comment '订单号',
   item                 int comment '订单item',
   content              varchar(255) comment '资料内容',
   imgs                 varchar(255) comment '图片',
   status               varchar(255) comment '状态',
   time                 timestamp,
   add1                 varchar(255),
   add2                 varchar(255),
   primary key (tid)
);

alter table ODS_ORDER_REFUND comment '退款表';

/*==============================================================*/
/* Table: ODS_RECORD                                            */
/*==============================================================*/
create table ODS_RECORD
(
   tid                  int not null AUTO_INCREMENT,
   name                 varchar(255),
   type                 varchar(255),
   money                float,
   uid                  int,
   pay_type             varchar(255),
   target_id            varchar(255),
   target_type          varchar(255),
   ono                  varchar(255),
   time                 timestamp,
   status               varchar(255),
   add1                 varchar(255),
   add2                 varchar(255),
   primary key (tid)
);

alter table ODS_RECORD comment '用户支出收入记录';

INSERT INTO ods_order_env (tid, akey, aval, type, status, time, add1, add2)
VALUES
	(1, 'RCP', 'http://rcp.dev.jxzy.com/usr/purchase-course?', 'PAID_CB', 'N', '2015-12-21 15:00:42', NULL, NULL),
   (2, 'RCP_HDKX', 'http://www.hd.kuxiao.cn/usr/purchase-course?', 'PAID_CB', 'N', '2015-12-21 15:00:42', NULL, NULL),
   (3, 'RCP_KX', 'http://www.kuxiao.cn/usr/purchase-course?', 'PAID_CB', 'N', '2015-12-21 15:00:42', NULL, NULL),
   (4, 'RCP_CHK', 'http://rcp.chk.jxzy.com/usr/purchase-course?', 'PAID_CB', 'N', '2015-12-21 15:00:42', NULL, NULL);

`
