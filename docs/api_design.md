## what informations need to get?
- profile informations of user (user name, bird date, sex, email, avatar, ...)
- their posted posts (in case: when i stalk my friend account)
- thier followee, follower (in case: when i stalk my friend account)
- Notifications of new posts of their followers
- NewFeeds of their followers when scrolling through Facebook (apply Pagination and filters)
- number of 'like', type of like, who is like in the specific post.
- all comments of a post
- search by user name someone to stalk them
.....
## what informations need to write?
- register/delete account (user_name, password)
- loggin
- change password
- edit profile information
- create/edit/delete a post
- follow/unfollow someone
- like/unlike at any post
- comment at any post
- edit/delete comment
.....

## Design API interface
# Read
- profile informations of user (user name, bird date, sex, email, avatar, ...)
    - GET : users/user_id
- their posted posts (in case: when i stalk my friend account)
    - GET : posts/user_id
- thier followees, followers (in case: when i stalk my friend account)
    - GET : followers/user_id
    - GET : followees/user_id
- Notifications of new posts of their followers
    - GET : notifications/
- NewFeeds of their followers when scrolling through Facebook (apply Pagination and filters)
    - GET : feeds?offset=?limit=
- number of 'like', type of like, who is like in the specific post.
    - GET : reactions/user_id/post_id
- all comments of a post
    - GET : comments/user_id/post_id
- search by user name someone to stalk them
    - GET : search/user_name
# Write
- register/delete account (user_name, password)
    - POST  : users/ {user_name, password}
    - DELETE: users/user_id
- loggin:
    - POST: loggin/{user_name, password}
- change password:
    - PUT : loggin?password=
- edit profile information:
    - PUT: users/{user name, bird date, sex, email, avatar, ...}
- create/edit/delete a post:
    - POST   : posts/{content, image/video}
    - PUT    : posts/{content, image/video}
    - DELETE : posts/post_id
- follow/unfollow someone
    - POST   : followers/user_id
    - DELETE :followers/user_id
- like/unlike at any post
    - POST   : reactions/user_id/post_id/{type}
    - DELETE : reactions/user_id/post_id
- comment at any post
    - POST : comments/user_id/post_id/{content}
- edit/delete comment
    - PUT    : comments/user_id/post_id/comment_id/{content}
    - DELETE : comments/user_id/post_id/comment_id

##### ###################################################
# API Design Document
## Dựa trên thiết kế API mà bạn đã cung cấp và các nhận xét trước đó, tôi sẽ xây dựng một phiên bản cải tiến. Thiết kế này sẽ tuân thủ các tiêu chuẩn RESTful, đảm bảo tính nhất quán, bảo mật (với xác thực JWT), và khả năng mở rộng. Tôi cũng sẽ tích hợp các đề xuất như phân trang, phản hồi chi tiết, và hỗ trợ phân vùng (liên kết với SQL trước đó). Dưới đây là thiết kế API mới:

---

### Thiết kế API cải tiến

#### 1. Auth
- **Register**
  - `POST /auth/register`
  - **Body**: `{username: string, email: string, password: string}`
  - **Header**: Không yêu cầu
  - **Response**: 
    - `200 OK`: `{user_id: number, token: string}`
    - `400 Bad Request`: `{error: "Invalid data"}`
- **Login**
  - `POST /auth/login`
  - **Body**: `{username: string, password: string}`
  - **Header**: Không yêu cầu
  - **Response**: 
    - `200 OK`: `{token: string}`
    - `401 Unauthorized`: `{error: "Invalid credentials"}`
- **Change Password**
  - `PUT /users/{user_id}/password`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{old_password: string, new_password: string}`
  - **Response**: 
    - `200 OK`: `{message: "Password updated"}`
    - `403 Forbidden`: `{error: "Invalid old password"}`
- **Delete Account**
  - `DELETE /users/{user_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{message: "Account soft deleted"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`
  - **Note**: Cập nhật `isDeleted = true` trong bảng `User`.

#### 2. Profile
- **Get User Profile**
  - `GET /users/{user_id}`
  - **Header**: `Authorization: Bearer <token>` (tùy chọn, chỉ bạn bè hoặc công khai)
  - **Response**: 
    - `200 OK`: `{user_id: number, username: string, email: string, avatar: string, bio: string, createdAt: string, isDeleted: boolean}`
    - `403 Forbidden`: `{error: "Private profile"}`
- **Update User Profile**
  - `PUT /users/{user_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{avatar: string, username: string, bio: string}` (partial update)
  - **Response**: 
    - `200 OK`: `{message: "Profile updated"}`
    - `400 Bad Request`: `{error: "Invalid data"}`
- **Search Users**
  - `GET /users?search={query}&offset={number}&limit={number}&sort={field}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{users: [{user_id, username, avatar}, ...], total: number}`
    - `400 Bad Request`: `{error: "Invalid parameters"}`

#### 3. Posts
- **Get Post**
  - `GET /posts/{post_id}`
  - **Header**: `Authorization: Bearer <token>` (tùy chọn)
  - **Response**: 
    - `200 OK`: `{post_id: number, user_id: number, content: string, createdAt: string, isDeleted: boolean, media_ids: [number]}`
    - `404 Not Found`: `{error: "Post not found"}`
- **Get User Posts**
  - `GET /users/{user_id}/posts?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{posts: [{post_id, content, createdAt, isDeleted}, ...], total: number}`
    - `404 Not Found`: `{error: "User not found"}`
- **Create Post**
  - `POST /posts`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{content: string, media_ids: [number]}`
  - **Response**: 
    - `201 Created`: `{post_id: number, message: "Post created"}`
    - `400 Bad Request`: `{error: "Invalid data"}`
- **Update Post**
  - `PUT /posts/{post_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{content: string, media_ids: [number]}`
  - **Response**: 
    - `200 OK`: `{message: "Post updated"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`
- **Delete Post**
  - `DELETE /posts/{post_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{message: "Post soft deleted"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`
  - **Note**: Cập nhật `isDeleted = true`.

#### 4. Reactions
- **Get Likes**
  - `GET /posts/{post_id}/likes`
  - **Header**: `Authorization: Bearer <token>` (tùy chọn)
  - **Response**: 
    - `200 OK`: `{count: number, types: [string], users: [{user_id, username}, ...]}`
    - `404 Not Found`: `{error: "Post not found"}`
- **Like Post**
  - `POST /posts/{post_id}/likes`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{reaction_type: string}` (lấy `user_id` từ token)
  - **Response**: 
    - `201 Created`: `{message: "Liked"}`
    - `400 Bad Request`: `{error: "Invalid reaction type"}`
- **Unlike Post**
  - `DELETE /posts/{post_id}/likes`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{message: "Unliked"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`

#### 5. Comments
- **Get Comments**
  - `GET /posts/{post_id}/comments?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>` (tùy chọn)
  - **Response**: 
    - `200 OK`: `{comments: [{comment_id, user_id, content, createdAt, isDeleted}, ...], total: number}`
    - `404 Not Found`: `{error: "Post not found"}`
- **Create Comment**
  - `POST /posts/{post_id}/comments`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{content: string}` (lấy `user_id` từ token)
  - **Response**: 
    - `201 Created`: `{comment_id: number, message: "Comment created"}`
    - `400 Bad Request`: `{error: "Invalid content"}`
- **Update Comment**
  - `PUT /comments/{comment_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `{content: string}`
  - **Response**: 
    - `200 OK`: `{message: "Comment updated"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`
- **Delete Comment**
  - `DELETE /comments/{comment_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{message: "Comment soft deleted"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`
  - **Note**: Cập nhật `isDeleted = true`.

#### 6. Follows
- **Get My Followers**
  - `GET /me/followers?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>` (bắt buộc)
  - **Response**: 
    - `200 OK`: `{followers: [{user_id, username, avatar}, ...], total: number}`
    - `401 Unauthorized`: `{error: "Unauthorized"}`

- **Get My Following**
  - `GET /me/following?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>` (bắt buộc)
  - **Response**: 
    - `200 OK`: `{following: [{user_id, username, avatar}, ...], total: number}`
    - `401 Unauthorized`: `{error: "Unauthorized"}`
- **Get Followers**
  - `GET /users/{user_id}/followers?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>` (tùy chọn)
  - **Response**: 
    - `200 OK`: `{followers: [{user_id, username, avatar}, ...], total: number}`
    - `404 Not Found`: `{error: "User not found"}`
- **Get Following**
  - `GET /users/{user_id}/following?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>` (tùy chọn)
  - **Response**: 
    - `200 OK`: `{following: [{user_id, username, avatar}, ...], total: number}`
    - `404 Not Found`: `{error: "User not found"}`
- **Follow User**
  - `POST /users/{target_user_id}/follow`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `201 Created`: `{message: "Followed"}`
    - `400 Bad Request`: `{error: "Already following"}`
- **Unfollow User**
  - `DELETE /users/{target_user_id}/follow`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**: 
    - `200 OK`: `{message: "Unfollowed"}`
    - `403 Forbidden`: `{error: "Unauthorized"}`

#### 7. Feeds & Notifications
- **Get My News Feed**
  - `GET /feeds?before={timestamp}&limit={number}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**:
    - `200 OK`: {
        feeds: [{
          post_id, 
          user_id, 
          username, 
          avatar, 
          content, 
          media_urls, 
          created_at,
          like_count,
          comment_count,
          is_liked: true/false
        }, ...],
        next_cursor: timestamp (or post_id)
      }
    - `401 Unauthorized`: `{error: "Unauthorized"}`
- **Get Notifications**
  - `GET /notifications?offset={number}&limit={number}`
  - **Header**: `Authorization: Bearer <token>`
  - **Response**:
    - `200 OK`: `[{id, type, source_user_id, post_id, read, created_at}]`
    - `400 Bad Request`: `{error: "Invalid parameters"}`
- **Mark Notification as Read**
  - `PATCH /notifications/{notification_id}`
  - **Header**: `Authorization: Bearer <token>`
  - **Request Body** (optional): `{ read: true }` (hoặc không cần nếu mặc định là đánh dấu đã đọc)
  - **Response**:
    - `200 OK`: `{ message: "Notification marked as read" }`
    - `404 Not Found`: nếu không tìm thấy notification_id
    - `403 Forbidden`: nếu notification không thuộc về user đang đăng nhập

#### 8. Media
- **Upload Media**
  - `POST /media`
  - **Header**: `Authorization: Bearer <token>`
  - **Body**: `multipart/form-data` với `{type: string, file: file, post_id: number}`
  - **Response**: 
    - `201 Created`: `{media_id: number, message: "Media uploaded"}`
    - `400 Bad Request`: `{error: "Invalid data"}`
---