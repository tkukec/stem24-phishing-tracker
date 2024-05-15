### Technologies Used
1. **Golang**: Our server-side application is built using Go, a powerful and efficient programming language.
2. **JWT Tokens**: JSON Web Tokens (JWT) are used for authentication and authorization purposes, providing secure access to our server's resources.
3. **Roles on User Model**: Users are assigned specific roles, allowing for fine-grained access control and permission management.
4. **Password Hashing**: User passwords are securely hashed before storage, ensuring the confidentiality and integrity of user credentials.
5. **Metrics with Prometheus**: Prometheus is integrated into our server application for monitoring and collecting metrics, enabling better insights into performance and behavior.
6. **Tenants**: Our server supports multi-tenancy, allowing for the isolation and segregation of data between different clients or user groups.
7. **2FA with Google Auth**: Two-Factor Authentication (2FA) using Google Authenticator enhances the security of user accounts by requiring an additional verification step.
8. **Form Validation on Backend**: Backend form validation is implemented to ensure data integrity and prevent malicious input.

### Frontend Integration
Our server integrates seamlessly with frontend applications, providing secure and efficient communication. Here are some key points regarding frontend integration:

1. **JWT for Protected Routes**: Frontend routes are protected using JWT tokens, ensuring that only authenticated users can access restricted resources.
2. **Form Validation**: Frontend form validation is implemented to provide a smooth and error-free user experience.
3. **Security-focused Packages**: We use the most up-to-date packages with no known security vulnerabilities to mitigate risks and ensure a robust security posture.

## Bonus task
- Using eslint, husky and lint staged
