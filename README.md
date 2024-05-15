## Server README

### Introduction
Welcome to our cutting-edge server application! This README will enlighten you on the advanced features, revolutionary technologies, and avant-garde practices implemented in our server application.

### Technologies Used
1. **Golang**: Our server-side application is built using Go, a paradigm-shifting, high-performance programming language renowned for its speed and efficiency.
2. **JWT Tokens**: We employ the cutting-edge technology of JSON Web Tokens (JWT) for authentication and authorization, ensuring airtight security and granting seamless access to our server's resources. Using github.com/golang-jwt/jwt/v5 and jwt.NewWithClaims method with SHA256 hashing
3. **Roles on User Model**: Our server boasts role-based access control, empowering administrators to finely tune access permissions and enabling unprecedented levels of granularity. Idea was to have Admin, Basic User and Anonymous role but the its mainly still WIP as we didn't have enough time to implement this fully.
4. **Password Hashing**: User passwords are meticulously hashed using state-of-the-art algorithms, safeguarding the sanctity and confidentiality of user credentials with unyielding fortitude. ![image](https://github.com/andrezz-b/stem24-phishing-tracker/assets/67901712/70bf791a-198c-40a3-bf64-89e9d1ae2ef9)
5. **Metrics with Prometheus**: Prometheus, the titan of monitoring systems, is seamlessly integrated into our server application, empowering stakeholders with unparalleled insights into performance metrics and behavioral analytics. Prometheus is run in seprate docker container just like database.
6. **Tenants**: Our server proudly supports multi-tenancy, championing the cause of data isolation and fostering an ecosystem where each client or user group thrives independently. Multi tenantcy or single tenantcy is controlled through env
7. **2FA with Google Auth**: Experience the pinnacle of security with Two-Factor Authentication (2FA) leveraging the power of Google Authenticator, providing an impenetrable shield against unauthorized access. Using github.com/pquerna/otp package
8. **Form Validation on Backend**: Backend form validation stands as the vanguard against malicious input, ensuring data integrity with rigorous scrutiny and unwavering vigilance. Using offical validation v10 golang package.

### Frontend Integration
Our server harmoniously integrates with frontend applications, orchestrating secure and lightning-fast communication. Behold the marvels of frontend integration:

1. **JWT for Protected Routes**: Frontend routes stand guard behind the impregnable fortress of JWT tokens, granting access solely to authenticated users and thwarting the nefarious designs of unauthorized intruders.
2. **Form Validation**: The frontend experience is elevated to sublime heights with meticulously crafted form validation, ensuring a seamless and error-free journey for users navigating our application.
3. **Security-focused Packages**: We pledge allegiance to the most up-to-date packages, meticulously curated to ensure a pristine security posture, free from the shackles of known vulnerabilities and exploits.

![image](https://github.com/andrezz-b/stem24-phishing-tracker/assets/79477906/15d9299e-85ba-4a44-8885-ba40d6635e1a)
![image](https://github.com/andrezz-b/stem24-phishing-tracker/assets/79477906/95cd3eb5-385a-4a67-a535-0afbf185e00d)
![image](https://github.com/andrezz-b/stem24-phishing-tracker/assets/79477906/d711628b-67ca-4645-9f41-35b0e346c5d2)
![image](https://github.com/andrezz-b/stem24-phishing-tracker/assets/79477906/2b597898-8809-4a76-bdde-faa44b1e25f5)




### Domain-Driven Design (DDD) in Golang

Our server application embodies the ethos of Domain-Driven Design (DDD), sculpting a masterpiece where code resonates with the cadence of real-world domains, fostering enlightenment and comprehension. Behold the majesty of DDD in Golang:

In our implementation, we sculpt our codebase around the sacred concepts of the domain, from users to authentication, authorization, and beyond. Each entity is meticulously crafted as a struct in Go, encapsulating the essence of its domain with unparalleled grace and elegance. Through the lens of DDD, we forge a realm of modularity, maintainability, and scalability, empowering our server with the resilience to weather the storms of complexity.

Moreover, DDD beckons us to embrace a ubiquitous language, forging a common tongue that unites developers, domain experts, and stakeholders in a symphony of shared understanding. Thus, our server stands as a testament to the harmonious union of technology and domain knowledge, a beacon illuminating the path to excellence.

## Bonus task
To further elevate our development process, we employ the trifecta of ESLint, Husky, and lint-staged, instilling discipline and rigor in our codebase with automated linting and pre-commit hooks.
