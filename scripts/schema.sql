-- Log initialization start
SELECT 'Initializing database schema...' as '';

use findata;

CREATE TABLE IF NOT EXISTS `ticks`
(
    `timestamp` bigint unsigned NOT NULL,
    `symbol`    varchar(8) NOT NULL,
    `bid`       float     NOT NULL,
    `ask`       float     NOT NULL,
    CONSTRAINT ticks_pk
        PRIMARY KEY (`timestamp`, `symbol`)
    -- # TODO maybe need some indexes ?
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `orders`
(
    `timestamp` bigint unsigned NOT NULL,
    `product_id`    varchar(8) NOT NULL,
    `type`      varchar(8) NOT NULL,
    `order_id`  varchar(64) NULL,
    `funds`     float NULL, -- funds in USD
    `side`      varchar(8) NULL, -- buy or sell
    `size`      float NULL, -- size of order
    `price`     float NULL, -- price of order
    `order_type` varchar(8) NULL, -- market or limit
    `client_oid` varchar(64) NULL, -- client order id
    `sequence` bigint unsigned NOT NULL,
    `remaining_size` float NULL,
    `reason` varchar(64) NULL,
    `trade_id` BIGINT NULL,
    `maker_order_id` varchar(64) NULL,
    `taker_order_id` varchar(64) NULL,
    CONSTRAINT orders_pk
        PRIMARY KEY (`timestamp`, `product_id`, `type`, `sequence`)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `swap_transactions`
(
    `tx_hash` varchar(66) NOT NULL,
    `version` varchar(8) NOT NULL,
    `exchange` varchar(100) NOT NULL, -- dex name (e.g. uniswap)
    `amount_in` varchar(100) NOT NULL,
    `to_address` varchar(42) NOT NULL,
    -- `deadline` datetime NOT NULL,
    `token_path_from` varchar(42) NOT NULL,
    `token_path_to` varchar(42) NOT NULL,
    `value` float NOT NULL DEFAULT 0,
    `amount_token_desired` varchar(100) NULL, -- Uniswap V2 add liquidity
    `amount_token_min` varchar(100) NULL, -- Uniswap V2 add liquidity
    `amount_eth_min` varchar(100) NULL, -- Uniswap V2 add liquidity
    `method_id` varchar(10) NULL,
    `method_name` varchar(100) NULL,
    `liquidity` varchar(100) NULL, -- Uniswap V2 remove liquidity
    `last_updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT swap_transactions_pk
        PRIMARY KEY (`tx_hash`)
) ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `token_metadata` (
    `address` varchar(42) NOT NULL,
    `decimals` tinyint UNSIGNED NULL,
    `symbol` varchar(10) NULL,
    `price` float NULL,
    `last_updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`address`)
) ENGINE = InnoDB;
-- ALTER TABLE token_metadata ADD COLUMN `last_updated` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- ALTER TABLE token_metadata DROP COLUMN `last_updated`;