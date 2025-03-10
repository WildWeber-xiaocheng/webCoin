CREATE TABLE exchange_coin  (
                                `id` bigint(0) NOT NULL AUTO_INCREMENT,
                                `symbol` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '交易币种名称，格式：BTC/USDT',
                                `base_coin_scale` int(0) NOT NULL COMMENT '基币小数精度',
                                `base_symbol` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '结算币种符号，如USDT',
                                `coin_scale` int(0) NOT NULL COMMENT '交易币小数精度',
                                `coin_symbol` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '交易币种符号',
                                `enable` int(0) NOT NULL COMMENT '状态，1：启用，2：禁止',
                                `fee` decimal(8, 4) NOT NULL COMMENT '交易手续费',
                                `sort` int(0) NOT NULL COMMENT '排序，从小到大',
                                `enable_market_buy` int(0) NOT NULL DEFAULT 1 COMMENT '是否启用市价买',
                                `enable_market_sell` int(0) NOT NULL DEFAULT 1 COMMENT '是否启用市价卖',
                                `min_sell_price` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT '最低挂单卖价',
                                `flag` int(0) NOT NULL DEFAULT 0 COMMENT '标签位，用于推荐，排序等,默认为0，1表示推荐',
                                `max_trading_order` int(0) NOT NULL DEFAULT 0 COMMENT '最大允许同时交易的订单数，0表示不限制',
                                `max_trading_time` int(0) NOT NULL DEFAULT 0 COMMENT '委托超时自动下架时间，单位为秒，0表示不过期',
                                `min_turnover` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT '最小挂单成交额',
                                `clear_time` bigint(0) NOT NULL DEFAULT 0 COMMENT '清盘时间',
                                `end_time` bigint(0) NOT NULL DEFAULT 0 COMMENT '结束时间',
                                `exchangeable` int(0) NOT NULL DEFAULT 1 COMMENT ' 是否可交易',
                                `max_buy_price` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT '最高买单价',
                                `max_volume` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT '最大下单量',
                                `min_volume` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT '最小下单量',
                                `publish_amount` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT ' 活动发行数量',
                                `publish_price` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT ' 分摊发行价格',
                                `publish_type` int(0) NOT NULL DEFAULT 1 COMMENT '发行活动类型 1:无活动,2:抢购发行,3:分摊发行',
                                `robot_type` int(0) NOT NULL DEFAULT 0 COMMENT '机器人类型',
                                `start_time` bigint(0) NOT NULL DEFAULT 0 COMMENT '开始时间',
                                `visible` int(0) NOT NULL DEFAULT 1 COMMENT ' 前台可见状态',
                                `zone` int(0) NOT NULL DEFAULT 0 COMMENT '交易区域',
                                PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;

INSERT INTO exchange_coin(`symbol`, `base_coin_scale`, `base_symbol`, `coin_scale`, `coin_symbol`, `enable`, `fee`, `sort`, `enable_market_buy`, `enable_market_sell`, `min_sell_price`, `flag`, `max_trading_order`, `max_trading_time`,`min_turnover`, `clear_time`, `end_time`, `exchangeable`, `max_buy_price`, `max_volume`, `min_volume`, `publish_amount`, `publish_price`, `publish_type`, `robot_type`, `start_time`, `visible`, `zone`) VALUES ('BTC/USDT', 2, 'USDT', 2, 'BTC', 1, 0.0010, 1, 1, 1, 0.00000000, 1, 0, 0, 0.00000000, 1640998800000, 1640998800000, 1, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 1, 0, 1640998800000, 1, 0);
INSERT INTO exchange_coin(`symbol`, `base_coin_scale`, `base_symbol`, `coin_scale`, `coin_symbol`, `enable`, `fee`, `sort`, `enable_market_buy`, `enable_market_sell`, `min_sell_price`, `flag`, `max_trading_order`, `max_trading_time`, `min_turnover`, `clear_time`, `end_time`, `exchangeable`, `max_buy_price`, `max_volume`, `min_volume`, `publish_amount`, `publish_price`, `publish_type`, `robot_type`, `start_time`, `visible`, `zone`) VALUES ('ETH/USDT', 2, 'USDT', 2, 'ETH', 1, 0.0010, 3, 1, 1, 0.00000000, 0, 0, 0,  0.00000000, 1640998800000, 1640998800000, 1, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 0.00000000, 1, 0, 1640998800000, 1, 0);