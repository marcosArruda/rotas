# Rota de Viagem #

O Desafio proposto é o clássico problema do viajante pobre que tem todo o tempo do mundo :). Basicamente ele precisa escolher a rota mais barata para chegar ao seu destino sem se preocupar com a quantidade de conexões.

## Problemas da especificação ##

A especificação possui alguns problemas que encontrei durante a leitura inicial do desafio. O primeiro deles é o próprio exemplo dado como explicação. Nele, o viajante pretende ir de **GRU** até **CDG** porém o exemplo informa a alternativa incorreta como a resposta correta. O outro problema consiste que o Nó CDG é erradamente escrito no exemplo da interface shell. Nela, ecreve-se **CGD** ao invés de **CDG**(Charles de Gaulle) representando o famoso aeroporto na França, batizado em virtude do antigo Presidente francês de mesmo nome.

Em um caso real, a primeira coisa que eu faria antes de começar a escrever qualquer código seria questionar essas inconsistências com a pessoa que escreveu a especificação. Esses erros podem ser tanto propositáis com uma intenção oculta de negócio, ou apenas um erro humano. Nos dois casos, é necessária a correção pois, sendo especificação, ela DEVE ESPECIFICAR e não ocultar a regra de negócio.

A terceira falta encontrada foi na ausência da informação do sentido das rotas inputadas. Por exemplo, **GRU -> BRC = $10** significa que **BRC -> GRU = $10**?. Uma vez que não foi explicitado esse caso (numa situação real, também teria o meu questionamento), considero no arquivo compilado que f(GRU -> BRC) = f(BRC -> GRU), ou seja, o custo de **origem -> destino** é igual ao custo de **destino -> origem**.

## Como executar ##
1. Faça o download do binário relativo ao seu sistema operacional (compilei para **Linux**, **Windows** e **MacOS**)
2. Inicie o sistema executando:

```shell
$ ./rotas input-routes.csv
```

3. O sistema irá ler o arquivo csv de input e logo em seguida perguntar-a por uma **origem-destino** a ser calculada.

```shell
please enter the route: XXX-XXX
```
4. Depois de inserir a **origem-destino**, o sistema irá informar a rota mais barata encontrada.

```shell
best route: XXX - XXX - XXX - XXX - XXX > $XX
please enter the route: ZZZ-ZZZ
best route: ZZZ - ZZZ > $ZZ
```

## Estrutura dos arquivos/pacotes ##
Implementei a solução usando **Golang**, porém poderia ter utilizado **Java**, **Kotlin**, **Scala**, **Groovy** e **Python** que são outras linguágens que tenho proeficiência avançada. Na estrutura **Golang** temos:
* Pacotes (que também foram encapsulados como módulos do Golang 1.14, portanto Módulos e Pacotes são as mesmas coisas nessa minha implementação)
  * **graph**: Contem todo o algoritmo de parse e calculo e pesquisa das rotas, bem como a implementação manual de o algoritmo base da estrutura de arestas de um Gráfo. Eu sei que existem milhares de bibliotecas de gráfo por aí. Preferi implementar a minha na mão pelo desafio proposto.
  * **money**: Contem a estrutural de dados e algoritmos usados para o cálculo financeiro dos custos de cada uma das rotas.
  * **rest**: Algoritmos e estrutura relativas à interface rest utilizada.
* Third Party Libs:
  * **github.com/gorilla/mux**: Lib um pouco mais performática para a interface rest http (única lib de terceiros utilizada).
* Arquivos:
  * **\*.go**: são os arquivos fonte do Golang.
  * **\*_test.go**: Arquivo de testes unitários de cada pacote.
  * **go.mod & go.sum**: Arquivos da estrutura de módulos do Golang.
  * **main.go**: Arquivo principal de inicialização do sistema.
  * **graph/graph.go**: Módulo da implementação de Gráfos. Contém todo o algorítmo utilizado nas buscas e estrutura de dados.
  * **money/money.go**: Módulo de Money.
  * **rest/rest.go**: Módulo da Interface Rest implementada.

## Decisões de Design ##
* Utilização de Multithread para controlar a interface shell e a interface rest Http separadamente. Basicamente Cada chamada Rest resulta em uma nova thread no servidor Http, mas o servidor Http por sua vez é iniciado em outra thread separada da thread de leitura e calculo dos algoritmos. Isso resulta em um sistema mais performático.
* Encapsulamento e Modularização da implementação de Grafos dentro do pacote/módulo **graph**, de maneira que fica bastante semples utilizar os mesmos algoritmos em outro sistema futuro, seguindo as melhores práticas do desenvolvimento Golang.
* A funcionalidade de IO no arquivo de input em disco(CSV) foi colocada também dentro do módulo **graph** uma vez que é apenas o grafo que realiza este IO.
* Utilizei uma lib open-source para o roteamento das requisições Http. Esta lib é mais performática do que a lib default do Golang.

## A API Rest ##
A API Rest desenvolvida possui basicamente dois(s) endpoints. São eles:
* POST,PUT "/route" - Recebe no body um json no seguinte formato:
```
input:
  {
    "from": "XXX",
    "to": "XXX",
    "cost": 12.25
  }
output
{
  "from": "XXX",
  "to": "XXX",
  "cost": 12.25
}
```

* GET "/route/{from}/{to}" - Retorna a rota mais baráda encontrada no seguinte formato:
```
output:
{
  "best": "XXX - XXX - XXX",
  "cost": 12.25
}
```

ps.: Nenhuma rota exige autenticação;  
ps2.: Existe um script shell chamado **add-route.sh** que faz a chamada via curl para facilitar :).

## Pré-requisitos para Compilar ##
* Golang 1.14+

Para compilar, utilize o shell chamado **build.sh** que adicionei :)
