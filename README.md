# Projet Weather Data Collection and Retrieval

## Description

This project implements a system for collecting and retrieving meteorological data from airport sensors.
Data collected includes temperature, atmospheric pressure and wind speed.

## Structure

The project is structured as follows:

- `api/` contains OpenAPI/Swagger specifications, JSON schema files, protocol definition files
- `cmd/` contains all main applications for this project
- `internal/` contains the internal code specific to the project
- `docs/` contains the documentation for the project
- `scripts/` contains the scripts for building, testing, and other operations
- `test/` contains the tests for the project

## Installation

### Requirements

- Go 1.21.4
- Make 3.81 or compatible

### Installation steps

1. Clone the repository:
    ```bash
    git clone https://github.com/jarhead-killgrave/imt-2023-go-project-ZEBIAN-KITSOUKOU-MILLION-LEBRAS-ACHEK.git
    ```

2. Copy the `.env.dev` file to `.env`, and edit it to match your configuration:
   ```bash
   cp .env.example .env
   ```
   or for Windows:
   ```bash
	copy .env.example .env
   ```

3. Initialize all necessary services:
   ```bash
   make init
   ```
   Alternatively, you can initialize each service without using `make`:
   ```bash
   docker compose up -d
   ```
   After initialization, you should be able to access the following services:
	- [http://influxdb.metrics.meteo-airport.localhost](http://influxdb.metrics.meteo-airport.localhost) *(InfluxDB)*

## Usage

### Run the project

## Authors

- [ZEBIAN Jana](https://github.com/JanaZebian)
- [MILLION Julien](https://github.com/AlphaOrOmega)
- [LEBRAS Gregoire](https://github.com/gregoireLeBras)
- [KITSOUKOU Manne Émile](https://github.com/jarhead-killgrave)
- [ACHEK Jamil](https://github.com/JamWare)
