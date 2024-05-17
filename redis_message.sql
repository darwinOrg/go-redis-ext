CREATE TABLE `redis_message`
(
    `id`         bigint      NOT NULL AUTO_INCREMENT,
    `stream_key` varchar(64) NOT NULL COMMENT '消息流key',
    `message_id` varchar(64) NOT NULL COMMENT '消息id',
    `expired_at` datetime(3) NOT NULL COMMENT '过期时间',
    `created_by` bigint      NOT NULL,
    `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP (3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_message_id_stream_key` (`message_id`, `stream_key`)
) ENGINE = InnoDB COMMENT ='redis消息'
;