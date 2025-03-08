CREATE TABLE `member_wallet`  (
                                  `id` bigint(0) NOT NULL AUTO_INCREMENT,
                                  `address` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '充值地址',
                                  `balance` decimal(18, 8) NOT NULL COMMENT '可用余额',
                                  `frozen_balance` decimal(18, 8) NOT NULL COMMENT '冻结余额',
                                  `release_balance` decimal(18, 8) NOT NULL COMMENT '待释放余额',
                                  `is_lock` int(0) NOT NULL DEFAULT 0 COMMENT '钱包不是锁定 0 否 1 是',
                                  `member_id` bigint(0) NOT NULL COMMENT '用户id',
                                  `version` int(0) NOT NULL COMMENT '版本',
                                  `coin_id` bigint(0) NOT NULL COMMENT '货币id',
                                  `to_released` decimal(18, 8) NOT NULL COMMENT '待释放总量',
                                  `coin_name` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '货币名称',
                                  PRIMARY KEY (`id`) USING BTREE,
                                  UNIQUE INDEX `UKm68bscpof0bpnxocxl4qdnvbe`(`member_id`, `coin_id`) USING BTREE,
                                  INDEX `FKf9tgbp9y9py8t9c5xj0lllcib`(`coin_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci ROW_FORMAT = DYNAMIC;