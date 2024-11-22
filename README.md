# Watch Explanation On Youtube
[![Watch the video](https://drive.google.com/uc?export=view&id=1A22F7WFFO5kWiZ5oVdZ2mFXGBPYOHdbI)](https://www.youtube.com/watch?v=mvTzupOIGl0)

# Database Diagram
![Database Diagram](https://github.com/chrisprojs/Secure-Auth-With-Golang-Chi/blob/main/documentation/betamart-database.png)
## 1. public.users Table
This table stores the basic details of users. Relevant fields for email verification include:
<ul>
<li>user_id: A unique identifier for each user (primary key).</li>
<li>email: The user's email address.</li>
<li>isemailverified: A boolean indicating whether the email address has been verified.</li>
<li>isphonenumberverified: A boolean for phone verification (not directly related to email verification here).</li>
</ul>
