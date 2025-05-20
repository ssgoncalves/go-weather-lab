# Weather API - Projeto de Estudo em Go
Este repositório contém um projeto de estudo que utiliza Go para consultar dados de clima baseados em um CEP. O projeto utiliza o Docker para containerizar a aplicação e Docker Compose para facilitar a execução dos serviços.


## Como Rodar o Projeto
1. Clonar o Repositório
   Primeiro, clone o repositório para sua máquina local:

```bash
git clone https://github.com/seu-usuario/weather-api.git
cd weather-api
```

2. Copiar o arquivo .env.example para .env
   Dentro do repositório, você verá um arquivo chamado .env.example. Esse arquivo contém as variáveis de ambiente necessárias para a execução do projeto. Copie esse arquivo para um novo arquivo chamado .env.
```bash
cp .env.example .env
```

3. Definir a Chave da API no .env
   Abra o arquivo .env com um editor de texto (como o VS Code, nano, vim, etc.), e defina a chave da API do serviço de clima na variável WEATHER_API_KEY:

4. Rodar os Contêineres com Docker Compose
   Com o arquivo .env configurado, agora você pode iniciar a aplicação com o Docker Compose. Este comando irá construir as imagens e rodar o contêiner conforme a configuração definida no arquivo docker-compose.yml.

```bash
docker-compose up
```

O serviço estará disponível em http://localhost:8080/zip-code/SEU-CEP-AQUI/weather (ou a porta definida no docker-compose.yml).


Este comando irá rodar os testes dentro do contêiner, garantindo que todos os serviços e funcionalidades estão funcionando corretamente.

## Estrutura do Projeto

**cmd/:** Contém o arquivo main.go para rodar a aplicação.

**internal/api/:** Implementação do controlador de rotas (API).

**internal/application/:** Lógica de negócios e serviços.

**internal/domain/:** Definição de modelos e erros personalizados.

**internal/infrastructure/:** Implementação de integrações com APIs externas (CEP e Clima).