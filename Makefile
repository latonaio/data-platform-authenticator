DB_USER_NAME="XXXXXXXX"
DB_USER_PASSWORD="XXXXXXXX"
SRC_DIR="/app/src"
CONF_PATH="$(SRC_DIR)/mysql.conf"
CONTAINER_NAME="sample-mysql"

convert-json:
	npx json5 migration_env.json5 | jq . > migration_env.json

.PHONY: local-run
local-run:
	GO_ENV=dev go run ./cmd/server/.

# mac の場合の migration tool のインストール
install-migrate:
	brew install golang-migrate

# make migrate-up number=1
# User データがない状態で実行する
.PHONY: migrate-up
migrate-up: convert-json
	migrate -path db/migrations -database "$(shell jq -r .driver migration_env.json)://$(shell jq -r .user migration_env.json):$(shell jq -r .password migration_env.json)@tcp($(shell jq -r .address migration_env.json):$(shell jq -r .port migration_env.json))/$(shell jq -r .database migration_env.json)?multiStatements=true" up $(number)

.PHONY: migrate-force
migrate-force: convert-json
	migrate -path db/migrations -database "$(shell jq -r .driver migration_env.json)://$(shell jq -r .user migration_env.json):$(shell jq -r .password migration_env.json)@tcp($(shell jq -r .address migration_env.json):$(shell jq -r .port migration_env.json))/$(shell jq -r .database migration_env.json)?multiStatements=true" force $(number)

# make migrate-down number=1
.PHONY: migrate-down
migrate-down: convert-json
	migrate -path db/migrations -database "$(shell jq -r .driver migration_env.json)://$(shell jq -r .user migration_env.json):$(shell jq -r .password migration_env.json)@tcp($(shell jq -r .address migration_env.json):$(shell jq -r .port migration_env.json))/$(shell jq -r .database migration_env.json)?multiStatements=true" down $(number)

.PHONY: generate-key-pair
generate-key-pair:
	openssl genrsa 4096 > private.key
	openssl rsa -pubout < private.key > public.key

# ユーザーを作成する
create-user:
	docker exec -it $(CONTAINER_NAME) sh -c "mysql --defaults-extra-file=$(CONF_PATH) -t --show-warnings -e \"CREATE USER $(DB_USER_NAME)@localhost IDENTIFIED BY '$(DB_USER_PASSWORD)';\""

# ユーザー一覧を表示する
show-users:
	docker exec -it $(CONTAINER_NAME) sh -c "mysql --defaults-extra-file=$(CONF_PATH) -e \"SELECT host, user FROM mysql.user;\""

# ユーザーに全ての権限を付与する
grant-authority:
	docker exec -it $(CONTAINER_NAME) sh -c "mysql --defaults-extra-file=$(CONF_PATH) -e \"GRANT ALL ON *.* TO $(DB_USER_NAME)@localhost;\""

# データベースを作成する
create-database:
	docker exec -it $(CONTAINER_NAME) sh -c "mysql -u$(DB_USER_NAME) -p$(DB_USER_PASSWORD) -e \"CREATE DATABASE IF NOT EXISTS Authenticator DEFAULT CHARACTER SET UTF8;\""

# データベース一覧を表示する
show-databases:
	docker exec -it $(CONTAINER_NAME) sh -c "mysql -u$(DB_USER_NAME) -p$(DB_USER_PASSWORD) -e \"SHOW DATABASES;\""

# # テーブルを作成する
# create-table:
# 	docker exec -it $(CONTAINER_NAME) sh -c "mysql -u$(DB_USER_NAME) -p$(DB_USER_PASSWORD) -e \ 
# 	cat << 'EOF'
# 	\"USE Authenticator; \
# 	CREATE TABLE `Users` ( \
#     `id` int AUTO_INCREMENT, \
# 		`login_id` varchar(20) NOT NULL DEFAULT '', \
# 		`password` varchar(70) NOT NULL DEFAULT '', \
# 		`qos` varchar(20) NOT NULL DEFAULT 'default', \
# 		`is_encrypt` BOOLEAN NOT NULL DEFAULT TRUE, \
# 		`last_login_at` datetime DEFAULT NULL, \
# 		`created_at` datetime DEFAULT NULL, \
# 		`updated_at` datetime DEFAULT NULL, \
# 		`deleted_at` datetime DEFAULT NULL, \
# 		UNIQUE KEY (`login_id`), \
# 		PRIMARY KEY (`id`) \
# 	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\""
# 	EOF

# ユーザーを削除する
delete-user:
	docker exec -it $(CONTAINER_NAME) sh -c "mysql --defaults-extra-file=$(CONF_PATH) -e \"DROP USER $(DB_USER_NAME)@localhost;\""