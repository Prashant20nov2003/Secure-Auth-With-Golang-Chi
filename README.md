# Watch Explanation On Youtube
[![Watch the video](https://drive.google.com/uc?export=view&id=1A22F7WFFO5kWiZ5oVdZ2mFXGBPYOHdbI)](https://www.youtube.com/watch?v=mvTzupOIGl0)

# Database Diagram
![Database Diagram](https://github.com/chrisprojs/Secure-Auth-With-Golang-Chi/blob/main/documentation/betamart-database.png)
## 1. public.users Table
<ul>
  <li>user_id: A unique identifier for each user (primary key).</li>
  <li>username: The user's username.</li>
  <li>email: The user's email address.</li>
  <li>phonenumber: The user's phonenumber.</li>
  <li>lastlogin: When the user login to account.</li>
  <li>isemailverified: A boolean indicating whether the email address has been verified.</li>
  <li>isphonenumberverified: A boolean for phone verification (not directly related to email verification here).</li>
</ul>

## 2. public.useremailverifications Table
<ul>
  <li>emailverify_id: A unique identifier for each email verification record (primary key).</li>
  <li>user_id: A foreign key referencing the user_id in the public.users table, linking the verification attempt to a specific user.</li>
  <li>expires_at: A timestamp specifying when the verification code will expire.</li>
  <li>verif_code: The email verification code (e.g., a 6-digit code).</li>
  <li>attempts: A counter for tracking how many times the user has attempted verification (e.g., max 3 attempts).</li>
  <li>used_for: Specifies the purpose of verification (e.g., "Verify Email" or "Forget Password").</li>
  <li>is_verified: A boolean indicating whether the verification was successful.</li>
</ul>

## 3. public.products Table
<ul>
  <li>product_id: A unique identifier for each product (primary key). This value is generated dynamically (e.g., using concat('PR-', ...)).</li>
  <li>user_id: A foreign key referencing the user_id field in the public.users table. This links the product to the user who created it.</li>
  <li>product_name: The name of the product.</li>
  <li>product_photo: A field that stores the path or identifier for the product's image.</li>
  <li>price: The price of the product, stored as an integer (e.g., in cents or the smallest currency unit).</li>
  <li>visibility: A boolean indicating whether the product is visible to other users or is hidden.</li>
</ul>
