# Go Stress Test

Aplicação CLI utilizada para realizar testes de carga, fornecendo um relatório de tempo de execução e parâmetros de sucesso e erro das requisições realizadas.

## Instruções de execução
Para a execução, recomenda-se a utilização do `docker`. 

Primeiro, execute o comando 

`docker build -t stress-test .` 

na raiz do projeto para a geração da imagem.

Após isto, utilize o comando 

`docker run stress-test --url=https://example.com, --requests=100 --concurrency=10`

, sendo as flags utlizadas:
- url (string): URL alvo para o teste de carga;
- requests (int): Número de requisições a serem enviadas para a realização do teste;
- concurrency (int): Quantas requisições serão enviadas paralelamente.

Os três parâmetors utilizados para os testes são obrigatórios.