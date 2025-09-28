package load

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Executa o stress test em uma dada URL.",
	Long: `O comando 'load' envia um número de requisições definido para uma
dada URL, com um nível especificado de concorrência. Por fim, este comando 
retorna um relatório dos resultados do teste.`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Erro: a flag --url é obrigatória.")
			os.Exit(1)
		}
		if requests <= 0 || concurrency <= 0 {
			fmt.Println("Erro: as flags --requests e --concurrency são obrigatórias e devem receber valores maiores que zero")
			os.Exit(1)
		}

		fmt.Printf("Iniciando o teste de carga...\n")
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("Número total de requisições: %d\n", requests)
		fmt.Printf("Número de requisições em paralelo: %d\n", concurrency)
		fmt.Printf("--------------------------------------\n")

		stresstester := NewTester(url, requests, concurrency)
		results := stresstester.Run()

		printReport(results)
	},
}

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "Um programa de teste de carga para URLs",
}

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.Flags().StringVarP(&url, "url", "u", "", "URL alvo")
	loadCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Número de requisições")
	loadCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 0, "Nível de concorrência")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printReport(report *Report) {
	fmt.Println("-------------------------------------------------")
	fmt.Println("           Relatório do Teste de Carga           ")
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Tempo total de Execução:\t\t%s\n", report.ExecutionTime.Round(time.Millisecond))
	fmt.Printf("Total de requisições enviadas:\t\t%d\n", report.TotalRequests)
	fmt.Printf("Total de requisições com sucesso:\t%d\n", report.SuccessfulRequests)
	fmt.Printf("Total de requisições falhas:\t\t%d\n", report.FailedRequests)

	if len(report.ErrorStatusDistribution) > 0 {
		fmt.Printf("\nDistribuição de status de erro:\n")
		for status, count := range report.ErrorStatusDistribution {
			fmt.Printf("%d: %d\n", status, count)
		}
	}
	fmt.Println("-------------------------------------------------")
}
