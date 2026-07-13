# Interface Contracts: Pi-hole v6 API

This contract documents the Pi-hole v6 REST API endpoints used by the synchronization tool.

## 1. Authentication
Pi-hole v6 REST API authenticates using the `sid` parameter in headers or cookies:
- Header: `X-FTL-SID: <sid_token>` or `sid: <sid_token>`

## 2. Get Custom DNS Hosts
List all currently configured custom DNS entries.

- **URL**: `GET /api/config/dns/hosts`
- **Headers**:
  - `Accept: application/json`
  - `X-FTL-SID: <token>`
- **Response Format (200 OK)**:
  ```json
  [
    "10.7.82.10 misty.fog.lodge.chalko.com",
    "10.7.82.100 grafana-101.fog.lodge.chalko.com"
  ]
  ```

## 3. Add Custom DNS Host
Add a single custom DNS host record.

- **URL**: `POST /api/config/dns/hosts`
- **Headers**:
  - `Content-Type: application/json`
  - `X-FTL-SID: <token>`
- **Request Body**:
  ```json
  "10.7.82.100 grafana-101.fog.lodge.chalko.com"
  ```
- **Response Format (200 OK)**: Empty body or success confirmation.

## 4. Delete Custom DNS Host
Delete a custom DNS host record.

- **URL**: `DELETE /api/config/dns/hosts/<url_encoded_ip_and_host>`
- **Headers**:
  - `X-FTL-SID: <token>`
- **Path Parameter**: `<ip>%20<hostname>` (e.g. `10.7.82.100%20grafana-101.fog.lodge.chalko.com`)
- **Response Format (200 OK)**: Empty body or success confirmation.
