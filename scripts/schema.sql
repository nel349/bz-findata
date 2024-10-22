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
    CONSTRAINT orders_pk
        PRIMARY KEY (`timestamp`, `product_id`, `type`, `sequence`)
) ENGINE = InnoDB;