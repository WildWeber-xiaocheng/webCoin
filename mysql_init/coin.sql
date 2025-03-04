CREATE TABLE `coin`  (
                         `id` int(0) NOT NULL AUTO_INCREMENT,
                         `name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '货币',
                         `can_auto_withdraw` int(0) NOT NULL COMMENT '是否能自动提币',
                         `can_recharge` int(0) NOT NULL COMMENT '是否能充币',
                         `can_transfer` int(0) NOT NULL COMMENT '是否能转账',
                         `can_withdraw` int(0) NOT NULL COMMENT '是否能提币',
                         `cny_rate` double NOT NULL COMMENT '对人民币汇率',
                         `enable_rpc` int(0) NOT NULL COMMENT '是否支持rpc接口',
                         `is_platform_coin` int(0) NOT NULL COMMENT '是否是平台币',
                         `max_tx_fee` double NOT NULL COMMENT '最大提币手续费',
                         `max_withdraw_amount` decimal(18, 8) NOT NULL COMMENT '最大提币数量',
                         `min_tx_fee` double NOT NULL COMMENT '最小提币手续费',
                         `min_withdraw_amount` decimal(18, 8) NOT NULL COMMENT '最小提币数量',
                         `name_cn` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '中文名称',
                         `sort` int(0) NOT NULL COMMENT '排序',
                         `status` tinyint(0) NOT NULL COMMENT '状态 0 正常 1非法',
                         `unit` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '单位',
                         `usd_rate` double NOT NULL COMMENT '对美元汇率',
                         `withdraw_threshold` decimal(18, 8) NOT NULL COMMENT '提现阈值',
                         `has_legal` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否是合法币种',
                         `cold_wallet_address` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '冷钱包地址',
                         `miner_fee` decimal(18, 8) NOT NULL DEFAULT 0.00000000 COMMENT '转账时付给矿工的手续费',
                         `withdraw_scale` int(0) NOT NULL DEFAULT 4 COMMENT '提币精度',
                         `account_type` int(0) NOT NULL DEFAULT 0 COMMENT '币种账户类型0：默认  1：EOS类型',
                         `deposit_address` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '充值地址',
                         `infolink` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '币种资料链接',
                         `information` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '币种简介',
                         `min_recharge_amount` decimal(18, 8) NOT NULL COMMENT '最小充值数量',
                         PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 9 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;


INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (1, 'Bitcoin', 0, 0, 1, 0, 0, 0, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '比特币', 1, 0, 'BTC', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (2, 'Bitcoincash', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '比特现金', 1, 0, 'BCH', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (3, 'DASH', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '达世币', 1, 0, 'DASH', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (4, 'Ethereum', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '以太坊', 1, 0, 'ETH', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (5, 'GalaxyChain', 1, 1, 1, 1, 1, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '银河链', 1, 0, 'GCC', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (6, 'Litecoin', 1, 0, 1, 1, 1, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '莱特币', 1, 0, 'LTC', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (7, 'SGD', 1, 1, 1, 1, 0, 1, 0, 0.0002, 500.00000000, 1, 1.00000000, '新币', 4, 0, 'SGD', 0, 0.10000000, 1, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);
INSERT INTO `mscoin`.`coin`(`id`, `name`, `can_auto_withdraw`, `can_recharge`, `can_transfer`, `can_withdraw`, `cny_rate`, `enable_rpc`, `is_platform_coin`, `max_tx_fee`, `max_withdraw_amount`, `min_tx_fee`, `min_withdraw_amount`, `name_cn`, `sort`, `status`, `unit`, `usd_rate`, `withdraw_threshold`, `has_legal`, `cold_wallet_address`, `miner_fee`, `withdraw_scale`, `account_type`, `deposit_address`, `infolink`, `information`, `min_recharge_amount`) VALUES (8, 'USDT', 1, 1, 1, 1, 0, 1, 0, 0.0002, 5.00000000, 0.0002, 0.00100000, '泰达币T', 1, 0, 'USDT', 0, 0.10000000, 0, '0', 0.00000000, 4, 0, '', '', '', 0.00000000);