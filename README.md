# Toll-Calculator

O projeto **Toll-Calculator** é uma aplicação distribuída baseada em microsserviços, desenvolvida para calcular pedágios de veículos em tempo real. Cada módulo tem uma responsabilidade específica, como receber dados, calcular distâncias ou atuar como um gateway. A comunicação entre os microsserviços ocorre por meio de **Kafka** para processamento assíncrono de mensagens e **sockets** para comunicação direta.

## Estrutura do Projeto

### Microsserviços

- **gate (gateway):**
  - Localização: `gateway/main.go`
  - Descrição: Serve como um gateway de entrada para a aplicação, recebendo e distribuindo as mensagens entre os outros microsserviços.
  - Comando de build: `@go build -o bin/gate gateway/main.go`

- **obu (On-Board Unit):**
  - Localização: `obu/main.go`
  - Descrição: Responsável por simular as unidades embarcadas nos veículos, coletando dados sobre a localização e transmitindo para outros serviços.
  - Comando de build: `@go build -o bin/obu obu/main.go`

- **receive (Receiver Service):**
  - Localização: `data_receive/*.go`
  - Descrição: Serviço que recebe os dados transmitidos pelas unidades OBU e os disponibiliza para processamento.
  - Comando de build: `@go build -o bin/receive data_receive/*.go`

- **calculator (Distance Calculator):**
  - Localização: `distance_calculator/*.go`
  - Descrição: Calcula as distâncias percorridas pelos veículos com base nos dados recebidos e calcula o pedágio correspondente.
  - Comando de build: `@go build -o bin/calculator distance_calculator/*.go`

- **agg (Aggregator):**
  - Localização: `aggregator/*.go`
  - Descrição: Agrega dados de várias fontes e realiza o processamento em lote para melhorar a eficiência do cálculo.
  - Comando de build: `@go build -o bin/agg aggregator/*.go`

### Arquivos de configuração e geração de código

- **proto:**
  - Localização: `types/ptypes.proto`
  - Descrição: Arquivo de definição do protocolo gRPC utilizado para comunicação entre os serviços.
  - Comando de build: `protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative types/ptypes.proto`

- **Prometheus:**
  - Localização: `.config/prometheus.yml`
  - Descrição: Configuração do Prometheus para monitoramento da aplicação.
  - Comando de execução: `@./bin/prometheus --config.file=.config/prometheus.yml`

## Comunicação Entre Microsserviços

### Kafka
- A comunicação assíncrona entre os serviços é realizada via **Apache Kafka**. Cada serviço é responsável por publicar e consumir mensagens em tópicos específicos, garantindo a escalabilidade e resiliência do sistema.

### Sockets
- A comunicação em tempo real entre alguns componentes utiliza **sockets**, permitindo a transmissão imediata de informações como atualizações de localização e comandos de controle.

## Pré-requisitos

- **Go 1.16+**: Linguagem utilizada no desenvolvimento do projeto.
- **Apache Kafka**: Sistema de mensageria para comunicação entre os microsserviços.
- **Prometheus**: Ferramenta de monitoramento para observar métricas do sistema.
- **Docker** (opcional): Para facilitar a configuração dos ambientes de desenvolvimento e produção.

## Como Executar

1. Faça o build dos microsserviços:
    ```bash
    make gate
    make obu
    make receive
    make calculator
    make agg
    ```

2. Suba os serviços do Kafka e Prometheus utilizando `docker-compose` (caso necessário):
    ```bash
    docker-compose up -d
    ```

3. Execute os serviços em binários:
    ```bash
    ./bin/gate
    ./bin/obu
    ./bin/receive
    ./bin/calculator
    ./bin/agg
    ```

## Monitoramento

Para monitorar o sistema, o Prometheus está configurado para capturar métricas dos microsserviços. Utilize a interface web do Prometheus para acessar as métricas:

- URL: [http://localhost:9090](http://localhost:9090)

## Contribuições

Sinta-se à vontade para contribuir com melhorias no código, correções de bugs ou novas funcionalidades. Para isso, faça um fork do projeto, crie uma branch com suas alterações e envie um pull request.
