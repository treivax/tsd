# 🔒 Sécurité et Authentification

**Date** : 2025-12-16  
**Version** : 1.0.0  
**Statut** : Documentation officielle

---

## Architecture Sécurité

```mermaid
graph TB
    subgraph "Client Layer"
        USER[👤 User]
        CLIENT[📡 TSD Client]
    end

    subgraph "Transport Security"
        TLS[🔐 TLS 1.3<br/>Encryption]
        CERT[📜 Certificates<br/>Server + Client]
    end

    subgraph "Authentication Layer"
        AUTH[🔑 Auth Manager]
        AUTHKEY[🔑 Auth Key]
        JWT[🎫 JWT Token]
    end

    subgraph "Authorization Layer"
        VALID[✅ Token Validator]
        VERIFY[🔍 Signature Verify]
        EXPIRE[⏰ Expiration Check]
    end

    subgraph "Server Layer"
        SERVER[🖥️ TSD Server]
        HANDLER[📋 Request Handler]
    end

    subgraph "Protected Resources"
        EXECUTE[⚙️ Execute Endpoint]
        METRICS[📊 Metrics]
        HEALTH[💚 Health]
    end

    USER -->|Request| CLIENT
    CLIENT -->|HTTPS| TLS
    TLS -->|Encrypted| CERT
    CERT -->|Secure Channel| AUTH
    
    CLIENT -.->|Auth Key| AUTHKEY
    CLIENT -.->|JWT| JWT
    
    AUTH --> VALID
    VALID --> VERIFY
    VERIFY --> EXPIRE
    
    EXPIRE -->|Valid| SERVER
    EXPIRE -->|Invalid| REJECT[❌ Reject]
    
    SERVER --> HANDLER
    HANDLER --> EXECUTE
    HANDLER --> METRICS
    HANDLER --> HEALTH

    style TLS fill:#27AE60,color:#fff,stroke:#333,stroke-width:3px
    style AUTH fill:#E74C3C,color:#fff,stroke:#333,stroke-width:2px
    style VALID fill:#F39C12,color:#fff,stroke:#333,stroke-width:2px
    style REJECT fill:#C0392B,color:#fff,stroke:#333,stroke-width:3px
```

---

## Flux d'Authentification JWT

```mermaid
sequenceDiagram
    participant Admin
    participant AuthCmd
    participant AuthModule
    participant Client
    participant Server

    Note over Admin,AuthModule: Phase 1: Génération des Credentials

    Admin->>AuthCmd: tsd auth generate-key
    AuthCmd->>AuthModule: GenerateKey()
    AuthModule->>AuthModule: crypto/rand 32 bytes
    AuthModule->>AuthModule: base64 encode
    AuthModule-->>Admin: auth_key saved to .tsd_auth_key

    Admin->>AuthCmd: tsd auth generate-jwt --key=xxx
    AuthCmd->>AuthModule: GenerateJWT(key, claims)
    AuthModule->>AuthModule: Create JWT header
    AuthModule->>AuthModule: Create JWT payload<br/>{sub, exp, iat}
    AuthModule->>AuthModule: Sign with HMAC-SHA256
    AuthModule-->>Admin: jwt_token

    Note over Client,Server: Phase 2: Utilisation du Token

    Client->>Client: Load JWT from config
    Client->>Server: POST /execute<br/>Authorization: Bearer jwt_token

    Server->>Server: Extract token from header
    Server->>Server: ValidateToken(token)
    Server->>Server: Parse JWT parts
    Server->>Server: Verify signature
    Server->>Server: Check expiration
    
    alt Token Valid
        Server->>Server: Process request
        Server-->>Client: 200 OK + Results
    else Token Invalid
        Server-->>Client: 401 Unauthorized
    end
```

---

## Types d'Authentification

```mermaid
graph TB
    subgraph "Auth Types"
        NONE[🔓 none<br/>No Auth]
        KEY[🔑 key<br/>Auth Key]
        JWT[🎫 jwt<br/>JSON Web Token]
    end

    subgraph "Use Cases"
        DEV[Development<br/>Testing]
        PROD[Production<br/>Simple]
        ENTERPRISE[Enterprise<br/>Advanced]
    end

    subgraph "Features"
        F1[✅ Quick Setup]
        F2[✅ Stateless]
        F3[✅ Expiration]
        F4[✅ Claims]
        F5[✅ Rotation]
    end

    NONE --> DEV
    KEY --> PROD
    JWT --> ENTERPRISE

    NONE -.-> F1
    KEY -.-> F1
    KEY -.-> F2
    JWT -.-> F2
    JWT -.-> F3
    JWT -.-> F4
    JWT -.-> F5

    style NONE fill:#95A5A6
    style KEY fill:#F39C12
    style JWT fill:#27AE60,stroke:#333,stroke-width:2px
```

---

## Structure JWT

```mermaid
graph LR
    subgraph "JWT Token Structure"
        HEADER[📋 Header<br/>{alg: HS256, typ: JWT}]
        PAYLOAD[📦 Payload<br/>{sub, exp, iat, custom}]
        SIGNATURE[✍️ Signature<br/>HMAC-SHA256]
    end

    subgraph "Encoding"
        B64H[Base64 Header]
        B64P[Base64 Payload]
        SIGN[HMAC Sign]
    end

    subgraph "Final Token"
        TOKEN[header.payload.signature]
    end

    HEADER --> B64H
    PAYLOAD --> B64P
    B64H --> SIGN
    B64P --> SIGN
    SIGN --> SIGNATURE
    
    B64H --> TOKEN
    B64P --> TOKEN
    SIGNATURE --> TOKEN

    style HEADER fill:#E8F4F8
    style PAYLOAD fill:#FFE5CC
    style SIGNATURE fill:#D5F4E6
    style TOKEN fill:#C3F0CA,stroke:#333,stroke-width:2px
```

**Exemple de JWT :**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiJ1c2VyMTIzIiwiZXhwIjoxNzM0MzYwMDAwLCJpYXQiOjE3MzQzNTY0MDB9.
7xKw9Y2Vp8qN3mF5tR1sL9jH4aU6bD8eK0vC2wX3gT5
```

---

## Validation de Token

```mermaid
stateDiagram-v2
    [*] --> ExtractToken: Request arrives
    
    ExtractToken --> CheckFormat: Get Authorization header
    CheckFormat --> TokenMissing: No token
    CheckFormat --> ParseToken: Token present
    
    TokenMissing --> Reject401: 401 Unauthorized
    
    ParseToken --> SplitParts: Split by '.'
    SplitParts --> InvalidFormat: Not 3 parts
    SplitParts --> DecodeHeader: Valid format
    
    InvalidFormat --> Reject401
    
    DecodeHeader --> DecodePayload
    DecodePayload --> VerifySignature
    
    VerifySignature --> SignatureInvalid: Wrong signature
    VerifySignature --> CheckExpiration: Signature valid
    
    SignatureInvalid --> Reject401
    
    CheckExpiration --> TokenExpired: exp < now
    CheckExpiration --> TokenValid: exp >= now
    
    TokenExpired --> Reject401
    TokenValid --> [*]: Allow access
    
    Reject401 --> [*]

    note right of VerifySignature
        Constant-time comparison
        Protection contre timing attacks
    end note
```

---

## Configuration TLS

```mermaid
graph TB
    subgraph "TLS Configuration"
        CONFIG[🔧 TLS Config]
    end

    subgraph "Certificate Management"
        GENCERT[🔐 Generate Certs]
        LOADCERT[📂 Load Certs]
        VERIFYCERT[✅ Verify Certs]
    end

    subgraph "Security Settings"
        TLS13[TLS 1.3 only]
        CIPHER[Strong Ciphers]
        VERIFY[Verify Client<br/>Optional]
    end

    subgraph "Usage"
        SERVER[Server Config]
        CLIENT[Client Config]
    end

    CONFIG --> GENCERT
    CONFIG --> LOADCERT
    CONFIG --> VERIFYCERT
    
    CONFIG --> TLS13
    CONFIG --> CIPHER
    CONFIG --> VERIFY
    
    TLS13 --> SERVER
    TLS13 --> CLIENT
    CIPHER --> SERVER
    CIPHER --> CLIENT
    VERIFY --> SERVER

    style CONFIG fill:#27AE60,color:#fff
    style TLS13 fill:#E74C3C,color:#fff
    style CIPHER fill:#F39C12,color:#fff
```

**Configuration par défaut :**
```go
MinVersion: tls.VersionTLS13
CipherSuites: [
    TLS_AES_128_GCM_SHA256
    TLS_AES_256_GCM_SHA384
    TLS_CHACHA20_POLY1305_SHA256
]
```

---

## Génération de Certificats

```mermaid
sequenceDiagram
    participant User
    participant AuthCmd
    participant TLSConfig
    participant Crypto
    participant FileSystem

    User->>AuthCmd: tsd auth generate-certs
    AuthCmd->>TLSConfig: GenerateSelfSignedCert()
    
    TLSConfig->>Crypto: Generate RSA 2048 key
    Crypto-->>TLSConfig: Private Key
    
    TLSConfig->>TLSConfig: Create certificate template
    Note over TLSConfig: Subject: CN=TSD Server<br/>ValidFor: 365 days<br/>KeyUsage: Digital Signature<br/>ExtKeyUsage: Server Auth
    
    TLSConfig->>Crypto: Sign certificate
    Crypto-->>TLSConfig: Certificate
    
    TLSConfig->>FileSystem: Write server.key
    TLSConfig->>FileSystem: Write server.crt
    TLSConfig-->>AuthCmd: Success
    
    AuthCmd-->>User: Certificates generated
```

---

## Endpoints Sécurisés

```mermaid
graph TB
    subgraph "Public Endpoints"
        HEALTH[GET /health<br/>❌ No Auth]
    end

    subgraph "Protected Endpoints"
        EXECUTE[POST /execute<br/>✅ Auth Required]
        METRICS[GET /metrics<br/>✅ Auth Required]
    end

    subgraph "Middleware Chain"
        M1[1. TLS Termination]
        M2[2. CORS Headers]
        M3[3. Auth Validation]
        M4[4. Request Logging]
    end

    subgraph "Handler"
        H[Request Handler]
    end

    HEALTH -->|Direct| H
    
    EXECUTE --> M1
    METRICS --> M1
    
    M1 --> M2
    M2 --> M3
    M3 --> M4
    M4 --> H

    style HEALTH fill:#95A5A6
    style EXECUTE fill:#E74C3C,color:#fff
    style METRICS fill:#E74C3C,color:#fff
    style M3 fill:#F39C12,stroke:#333,stroke-width:2px
```

---

## Bonnes Pratiques Sécurité

```mermaid
mindmap
    root((🔒 Security))
        Authentication
            JWT avec expiration
            Clés fortes 256 bits
            Rotation régulière
            Pas de credentials hardcodés
        Transport
            TLS 1.3 obligatoire
            Certificats valides
            Strong ciphers only
            Perfect Forward Secrecy
        Validation
            Constant-time comparison
            Input sanitization
            Type checking
            Boundary checks
        Monitoring
            Failed auth attempts
            Token expiration logs
            TLS handshake errors
            Rate limiting
```

---

## Threat Model

```mermaid
graph TB
    subgraph "Threats"
        T1[🎭 Man-in-the-Middle]
        T2[🔓 Credential Theft]
        T3[⏰ Replay Attacks]
        T4[🔍 Token Leakage]
    end

    subgraph "Mitigations"
        M1[✅ TLS 1.3 Encryption]
        M2[✅ Secure Storage]
        M3[✅ Token Expiration]
        M4[✅ HTTPS Only]
    end

    subgraph "Detection"
        D1[📊 Failed Auth Logs]
        D2[🚨 Anomaly Detection]
        D3[📈 Rate Monitoring]
    end

    T1 --> M1
    T2 --> M2
    T3 --> M3
    T4 --> M4
    
    M1 --> D1
    M2 --> D2
    M3 --> D3
    M4 --> D1

    style T1 fill:#C0392B,color:#fff
    style T2 fill:#C0392B,color:#fff
    style T3 fill:#C0392B,color:#fff
    style T4 fill:#C0392B,color:#fff
    style M1 fill:#27AE60,color:#fff
    style M2 fill:#27AE60,color:#fff
    style M3 fill:#27AE60,color:#fff
    style M4 fill:#27AE60,color:#fff
```

---

## Configuration Serveur Sécurisé

```yaml
# Exemple de configuration serveur TSD
server:
  address: "0.0.0.0:8443"
  
  tls:
    enabled: true
    cert_file: "certs/server.crt"
    key_file: "certs/server.key"
    min_version: "TLS13"
    
  auth:
    type: "jwt"
    secret_file: ".tsd_auth_key"
    
  timeouts:
    read: 15s
    write: 15s
    idle: 60s
    
  rate_limit:
    requests_per_second: 100
    burst: 200
```

---

## Audit Log

```mermaid
graph LR
    subgraph "Events"
        E1[🔑 Auth Attempt]
        E2[✅ Auth Success]
        E3[❌ Auth Failure]
        E4[📝 Request Executed]
        E5[🚫 Access Denied]
    end

    subgraph "Log Fields"
        F1[Timestamp]
        F2[Client IP]
        F3[User ID]
        F4[Action]
        F5[Result]
        F6[Duration]
    end

    subgraph "Storage"
        S1[📄 File Log]
        S2[📊 Metrics]
        S3[🔍 SIEM]
    end

    E1 --> F1
    E2 --> F1
    E3 --> F1
    E4 --> F1
    E5 --> F1
    
    F1 --> S1
    F2 --> S1
    F3 --> S1
    F4 --> S1
    F5 --> S1
    F6 --> S2
    
    S1 --> S3
    S2 --> S3

    style E3 fill:#E74C3C,color:#fff
    style E5 fill:#E74C3C,color:#fff
    style S3 fill:#3498DB,color:#fff
```

---

## Références

- [Architecture Globale](01-global-architecture.md)
- [Module Auth](../../auth/)
- [Module TLS Config](../../internal/tlsconfig/)
- [Documentation Sécurité](../../SECURITY.md)

---

**Maintenu par** : TSD Contributors  
**Dernière mise à jour** : 2025-12-16
