# Hướng dẫn thiết lập CI/CD

Hệ thống CI/CD này tự động deploy code lên máy local của bạn khi có push vào branch `main`, với khả năng tự động rollback nếu deploy thất bại và quản lý version.

## Tính năng

- ✅ Tự động test và lint code trước khi deploy
- ✅ Tự động build Docker images với version tagging
- ✅ Tự động deploy trên self-hosted runner (không cần SSH)
- ✅ Quản lý version tự động (từ git tag hoặc tự generate)
- ✅ Tự động backup trước khi deploy
- ✅ Tự động rollback nếu deploy thất bại
- ✅ Ghi log tất cả các bước deploy và lỗi
- ✅ Không deploy database migrations (theo yêu cầu)
- ✅ Tự động tạo git tag sau khi deploy thành công

## Yêu cầu

1. **GitHub Repository**: Dự án phải được push lên GitHub
2. **Self-hosted Runner**: Máy local phải có:
   - GitHub Actions self-hosted runner đã được cấu hình
   - Docker và Docker Compose
   - Go và Node.js (hoặc sẽ được cài tự động bởi actions)
   - Quyền truy cập vào thư mục dự án

## Cấu hình Self-hosted Runner

Nếu chưa có self-hosted runner, làm theo các bước sau:

1. Vào **Settings** → **Actions** → **Runners** → **New self-hosted runner**
2. Chọn OS và architecture của máy bạn
3. Chạy các lệnh được cung cấp trên máy local
4. Đảm bảo runner đang chạy: `sudo ./svc.sh status`

**Lưu ý:** Runner phải có quyền:
- Truy cập vào thư mục dự án
- Chạy Docker (thêm user vào docker group: `sudo usermod -aG docker $USER`)
- Pull/push code từ GitHub (nếu cần tạo tag)

## Quản lý Version

Hệ thống sử dụng **semantic versioning** (x.y.z) được lưu trong file `deploy/VERSION`.

### Tự động bump version dựa trên commit message:

Version sẽ tự động được tăng dựa trên từ khóa trong commit message:

- **"major"** → Tăng major version (1.0.0 → 2.0.0)
- **"minor"** → Tăng minor version (1.0.0 → 1.1.0)
- **"patch"** → Tăng patch version (1.0.0 → 1.0.1) - **mặc định nếu không có từ khóa**

**Ví dụ commit messages:**
```bash
git commit -m "Add new feature [minor]"
# Version: 1.0.0 → 1.1.0

git commit -m "Fix bug [patch]"
# Version: 1.1.0 → 1.1.1

git commit -m "Breaking change [major]"
# Version: 1.1.1 → 2.0.0

git commit -m "Regular commit"
# Version: 1.1.1 → 1.1.2 (mặc định patch)
```

**Lưu ý:** Từ khóa không phân biệt hoa thường và có thể nằm ở bất kỳ đâu trong commit message.

Version được:
- Đọc từ file `deploy/VERSION` khi deploy
- Tự động bump dựa trên commit message
- Tag Docker images với format `lexigo-{service}:{version}`
- Tạo git tag `v{version}` sau khi deploy thành công
- Lưu vào `deploy/versions/current-prod.txt`

## Cấu trúc thư mục

Sau khi setup, cấu trúc sẽ như sau:

```
language-coach/
├── .github/
│   └── workflows/
│       └── deploy.yml          # GitHub Actions workflow
├── scripts/
│   ├── deploy.sh                # Script deploy với rollback và versioning
│   ├── rollback.sh              # Script rollback
│   └── dev.sh                   # Script start development environment
├── deploy/
│   ├── VERSION                  # File chứa version hiện tại (x.y.z)
│   ├── logs/                    # Logs deployment (tự động tạo)
│   │   └── deploy.log
│   ├── backups/                 # Thư mục backup (tự động tạo)
│   │   └── prod-YYYYMMDD-HHMMSS/
│   │       ├── docker-compose.prod.yml
│   │       ├── version.txt
│   │       └── previous-images.txt
│   └── versions/                # Thư mục lưu version info
│       └── current-prod.txt
```

## Cách hoạt động

### 1. Khi có push vào `main`:

1. **Test & Lint**: Chạy tests và lint cho backend và frontend trên self-hosted runner
   - Go version: 1.25
   - Node.js version: 20.19
2. **Bump Version**: Tự động tăng version dựa trên commit message:
   - Tìm từ khóa "major", "minor", hoặc "patch" trong commit message
   - Bump version tương ứng (mặc định là patch nếu không tìm thấy)
   - Cập nhật file `deploy/VERSION`
3. **Deploy Production**:
   - Tạo backup của deployment hiện tại (kèm version info)
   - Build Docker images mới
   - Tag images với version: `lexigo-{service}:{version}`
   - Stop containers cũ
   - Start containers mới
   - Kiểm tra health của services
   - Lưu version info vào `deploy/versions/current-prod.txt`
   - Tạo git tag `v{version}` sau khi deploy thành công
   - Nếu thất bại → tự động rollback

### 2. Rollback tự động:

Nếu deploy thất bại, hệ thống sẽ:
1. Ghi log lỗi
2. Tự động chạy rollback script
3. Khôi phục containers về version trước đó (từ backup)
4. Upload logs và version info lên GitHub Actions artifacts

## Sử dụng thủ công

### Deploy thủ công:

```bash
# Deploy production
# Version sẽ tự động bump dựa trên commit message
./scripts/deploy.sh

# Hoặc với version cụ thể (bỏ qua auto-bump)
./scripts/deploy.sh 1.0.0
```

### Development:

```bash
# Start development environment
./scripts/dev.sh
```

### Rollback thủ công:

```bash
# Rollback về version mới nhất
./scripts/rollback.sh

# Rollback về version cụ thể (dùng backup timestamp)
./scripts/rollback.sh 20250127-143000
```

### Xem version hiện tại:

```bash
# Xem version đang chạy
cat deploy/versions/current-prod.txt

# Xem tất cả versions đã deploy
ls -la deploy/versions/

# Xem version trong backup
cat deploy/backups/prod-*/version.txt
```

### Xem logs:

```bash
# Xem log deploy
tail -f deploy/logs/deploy.log

# Xem log với grep
grep ERROR deploy/logs/deploy.log
grep SUCCESS deploy/logs/deploy.log
```

## Kiểm tra deployment

Sau khi deploy, kiểm tra:

```bash
# Xem trạng thái containers
cd deploy/compose
docker compose -f docker-compose.prod.yml ps

# Xem logs containers
docker compose -f docker-compose.prod.yml logs -f

# Kiểm tra health
docker compose -f docker-compose.prod.yml ps
```

## Troubleshooting

### Lỗi self-hosted runner:

1. Kiểm tra runner đang chạy:
   ```bash
   cd ~/actions-runner
   sudo ./svc.sh status
   ```

2. Xem logs runner:
   ```bash
   tail -f ~/actions-runner/_diag/Runner_*.log
   ```

3. Restart runner nếu cần:
   ```bash
   sudo ./svc.sh stop
   sudo ./svc.sh start
   ```

### Lỗi permission:

Đảm bảo user có quyền:
- Chạy Docker (thêm vào docker group):
  ```bash
  sudo usermod -aG docker $USER
  ```
- Ghi vào thư mục dự án
- Tạo thư mục `logs/` và `deploy/backups/`

### Lỗi containers không healthy:

1. Xem logs:
   ```bash
   docker compose -f deploy/compose/docker-compose.prod.yml logs
   ```

2. Kiểm tra environment variables:
   ```bash
   cat deploy/env/prod/backend.env
   cat deploy/env/prod/frontend.env
   ```

3. Test containers thủ công:
   ```bash
   docker compose -f deploy/compose/docker-compose.prod.yml up -d
   ```

## Backup và Version

Hệ thống tự động:
- Tạo backup trước mỗi lần deploy (kèm version info)
- Giữ lại 5 backup gần nhất
- Backup được lưu tại `deploy/backups/`
- Mỗi backup chứa:
  - `docker-compose.{env}.yml` - Cấu hình compose
  - `version.txt` - Thông tin version
  - `previous-images.txt` - Thông tin images cũ

Version được lưu tại:
- `deploy/versions/current-{env}.txt` - Version hiện tại
- `deploy/versions/{version}.txt` - Version cụ thể
- Git tags (nếu được tạo tự động)

## Logs

Tất cả logs được ghi vào:
- **Local**: `deploy/logs/deploy.log`
- **GitHub Actions**: Artifacts (tự động upload sau mỗi deployment)

Logs bao gồm:
- Timestamp
- Log level (INFO, SUCCESS, WARNING, ERROR)
- Chi tiết từng bước deploy
- Version được deploy
- Lỗi nếu có

## Lưu ý

- ⚠️ **Chỉ deploy production**: Script chỉ hỗ trợ deploy production environment
- ⚠️ **Version format**: Phải là semantic versioning (x.y.z) trong file `deploy/VERSION`
- ⚠️ **Auto version bump**: Version tự động tăng dựa trên commit message (major/minor/patch)
- ⚠️ **Không deploy database migrations**: Script không chạy migrations tự động
- ⚠️ **Backup tự động**: Luôn có backup trước khi deploy
- ⚠️ **Rollback tự động**: Nếu deploy fail, sẽ tự động rollback
- ⚠️ **Health check**: Script đợi 60 giây để services healthy
- ⚠️ **Image tagging**: Images được tag với format `lexigo-{service}:{version}`
