package matchmaking;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

import matchmaking.saida.Resultado;

/** [P4-ETAPA-04] Matchmaking no paradigma orientado a objetos. */
public final class Main {
    private Main() {
    }

    public static void main(String[] args) {
        try {
            if (args.length > 0 && args[0].equals("-relatorio")) {
                if (args.length > 2) {
                    falhaUso("uso: java matchmaking.Main -relatorio [pasta]");
                }
                Path pasta = args.length == 2 ? Path.of(args[1]) : Path.of("..", "testes");
                System.out.print(new Relatorio(new Aplicacao()).gerar(pasta));
                return;
            }
            if (args.length == 1 && Files.isDirectory(Path.of(args[0]))) {
                System.out.print(new Relatorio(new Aplicacao()).gerar(Path.of(args[0])));
                return;
            }
            if (args.length > 1) {
                falhaUso("uso: java matchmaking.Main [arquivo | -relatorio pasta]");
            }
            String texto = args.length == 1
                    ? Files.readString(Path.of(args[0]), StandardCharsets.UTF_8)
                    : new String(System.in.readAllBytes(), StandardCharsets.UTF_8);
            Resultado resultado = new Aplicacao().executar(texto);
            System.out.print(resultado.texto());
            System.exit(resultado.codigoDeSaida());
        } catch (IOException ex) {
            System.err.println("falha ao ler arquivo: " + ex.getMessage());
            System.exit(2);
        }
    }

    private static void falhaUso(String mensagem) {
        System.err.println(mensagem);
        System.exit(2);
    }
}
