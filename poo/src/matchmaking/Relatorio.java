package matchmaking;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Stream;

import matchmaking.saida.Resultado;

public final class Relatorio {
    private final Aplicacao aplicacao;

    public Relatorio(Aplicacao aplicacao) {
        this.aplicacao = aplicacao;
    }

    public String gerar(Path pasta) throws IOException {
        if (!Files.isDirectory(pasta)) {
            throw new IOException("nenhum caso em " + pasta);
        }
        List<Path> entradas = new ArrayList<Path>();
        try (Stream<Path> listagem = Files.list(pasta)) {
            listagem.filter(caminho -> caminho.getFileName().toString().endsWith(".txt"))
                    .forEach(entradas::add);
        }
        if (entradas.isEmpty()) {
            throw new IOException("nenhum caso em " + pasta);
        }
        entradas.sort(Comparator
                .comparingInt((Path caminho) -> faixa(caminho.getFileName().toString()))
                .thenComparing(caminho -> caminho.getFileName().toString()));

        StringBuilder texto = new StringBuilder();
        texto.append(String.format("%-6s %-18s %-24s %s%n", "caso", "status", "detalhe", "espera"));
        int partidas = 0;
        int nenhuma = 0;
        int erros = 0;
        for (Path entrada : entradas) {
            String conteudo = Files.readString(entrada, StandardCharsets.UTF_8);
            Resultado resultado = aplicacao.executar(conteudo);
            switch (resultado.status()) {
                case "PARTIDA_FORMADA":
                    partidas++;
                    break;
                case "NENHUMA_PARTIDA":
                    nenhuma++;
                    break;
                case "ERRO_ENTRADA":
                    erros++;
                    break;
                default:
                    break;
            }
            String nome = entrada.getFileName().toString().replaceFirst("\\.txt$", "");
            texto.append(String.format("%-6s %-18s %-24s %s%n",
                    nome, resultado.status(), resultado.detalhe(), resultado.textoEspera()));
        }
        texto.append(String.format("%n%d casos: %d partidas, %d sem partida, %d erros de entrada%n",
                entradas.size(), partidas, nenhuma, erros));
        // O contrato usa LF. No Windows, %n grava CRLF.
        return texto.toString().replace("\r\n", "\n");
    }

    private static int faixa(String nome) {
        if (nome.startsWith("n")) {
            return 0;
        }
        if (nome.startsWith("l")) {
            return 1;
        }
        if (nome.startsWith("i")) {
            return 2;
        }
        return 3;
    }
}
