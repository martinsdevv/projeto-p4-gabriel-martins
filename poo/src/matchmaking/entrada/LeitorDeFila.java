package matchmaking.entrada;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

import matchmaking.dominio.Fila;
import matchmaking.dominio.Inscricao;
import matchmaking.dominio.Posicao;
import matchmaking.saida.ErroEntrada;

/**
 * Transforma o texto do contrato em inscricoes.
 * Um registro lido ainda nao e jogador: so vira Jogador se todas as regras passarem.
 */
public final class LeitorDeFila {
    private final List<RegraDeEntrada> regras = List.of(
            new AtributosDoJogador(),
            new IdentificadoresUnicos(),
            new TamanhoDosGrupos());

    public Leitura ler(String texto) {
        List<Registro> registros = new ArrayList<Registro>();
        String[] linhas = texto.split("\n", -1);
        int numero = 0;
        for (String bruta : linhas) {
            String linha = bruta.trim();
            if (linha.isEmpty()) {
                continue;
            }
            numero++;
            Optional<ErroEntrada> erro = lerLinha(linha, numero, registros);
            if (erro.isPresent()) {
                return Leitura.rejeitada(erro.get());
            }
        }
        for (RegraDeEntrada regra : regras) {
            Optional<ErroEntrada> erro = regra.verificar(registros);
            if (erro.isPresent()) {
                return Leitura.rejeitada(erro.get());
            }
        }
        List<Inscricao> inscricoes = new ArrayList<Inscricao>();
        for (Registro registro : registros) {
            inscricoes.add(paraInscricao(registro));
        }
        return Leitura.ok(Fila.formar(inscricoes));
    }

    private Optional<ErroEntrada> lerLinha(String linha, int numero, List<Registro> registros) {
        String[] partes = linha.split("\\|", -1);
        if (partes.length != 8) {
            return Optional.of(new ErroEntrada(
                    "CAMPO_OBRIGATORIO",
                    "linha " + numero + ": campo obrigatorio ausente"));
        }
        for (int i = 0; i < partes.length; i++) {
            partes[i] = partes[i].trim();
        }
        String id = partes[0];
        if (id.isEmpty()) {
            return Optional.of(new ErroEntrada(
                    "CAMPO_OBRIGATORIO",
                    "linha " + numero + ": identificador ausente"));
        }
        Integer habilidade = inteiro(partes[1]);
        Integer espera = inteiro(partes[4]);
        Integer latencia = inteiro(partes[5]);
        if (habilidade == null || espera == null || latencia == null) {
            return Optional.of(new ErroEntrada(
                    "CAMPO_OBRIGATORIO",
                    "linha " + numero + ": campo numerico ilegivel (jogador " + id + ")"));
        }
        if (partes[7].isEmpty()) {
            return Optional.of(new ErroEntrada(
                    "CAMPO_OBRIGATORIO",
                    "linha " + numero + ": regiao ausente (jogador " + id + ")"));
        }
        registros.add(new Registro(
                id,
                habilidade.intValue(),
                partes[2],
                alternativas(partes[3]),
                espera.intValue(),
                latencia.intValue(),
                grupo(partes[6]),
                partes[7]));
        return Optional.empty();
    }

    private static Inscricao paraInscricao(Registro registro) {
        Posicao preferida = Posicao.reconhecer(registro.preferida).orElseThrow();
        List<Posicao> alternativas = new ArrayList<Posicao>();
        for (String nome : registro.alternativas) {
            alternativas.add(Posicao.reconhecer(nome).orElseThrow());
        }
        return new Inscricao(
                registro.id,
                registro.habilidade,
                preferida,
                alternativas,
                registro.espera,
                registro.latencia,
                registro.grupo,
                registro.regiao);
    }

    private static List<String> alternativas(String campo) {
        List<String> nomes = new ArrayList<String>();
        if (campo.isEmpty() || campo.equals("-")) {
            return nomes;
        }
        String[] pedacos = campo.split(",");
        for (String pedaco : pedacos) {
            String nome = pedaco.trim();
            if (nome.isEmpty() || nome.equals("-") || nomes.contains(nome)) {
                continue;
            }
            nomes.add(nome);
        }
        return nomes;
    }

    private static String grupo(String campo) {
        if (campo.isEmpty() || campo.equals("-")) {
            return null;
        }
        return campo;
    }

    private static Integer inteiro(String campo) {
        try {
            return Integer.valueOf(campo);
        } catch (NumberFormatException ex) {
            return null;
        }
    }

    private interface RegraDeEntrada {
        Optional<ErroEntrada> verificar(List<Registro> registros);
    }

    private static final class AtributosDoJogador implements RegraDeEntrada {
        @Override
        public Optional<ErroEntrada> verificar(List<Registro> registros) {
            for (Registro jogador : registros) {
                Optional<ErroEntrada> erro = verificarJogador(jogador);
                if (erro.isPresent()) {
                    return erro;
                }
            }
            return Optional.empty();
        }

        private Optional<ErroEntrada> verificarJogador(Registro jogador) {
            if (jogador.habilidade < 0) {
                return Optional.of(new ErroEntrada(
                        "HABILIDADE_INVALIDA",
                        "habilidade deve ser maior ou igual a 0 (jogador " + jogador.id + ")"));
            }
            if (jogador.espera < 0) {
                return Optional.of(new ErroEntrada(
                        "ESPERA_INVALIDA",
                        "espera_s deve ser maior ou igual a 0 (jogador " + jogador.id + ")"));
            }
            if (jogador.latencia < 0) {
                return Optional.of(new ErroEntrada(
                        "LATENCIA_INVALIDA",
                        "latencia_ms deve ser maior ou igual a 0 (jogador " + jogador.id + ")"));
            }
            if (jogador.preferida.isEmpty()) {
                return Optional.of(new ErroEntrada(
                        "PREFERIDA_AUSENTE",
                        "posicao preferida ausente (jogador " + jogador.id + ")"));
            }
            if (Posicao.reconhecer(jogador.preferida).isEmpty()) {
                return Optional.of(new ErroEntrada(
                        "POSICAO_DESCONHECIDA",
                        "posicao nao reconhecida: " + jogador.preferida + " (jogador " + jogador.id + ")"));
            }
            for (String alternativa : jogador.alternativas) {
                if (Posicao.reconhecer(alternativa).isEmpty()) {
                    return Optional.of(new ErroEntrada(
                            "POSICAO_DESCONHECIDA",
                            "posicao nao reconhecida: " + alternativa + " (jogador " + jogador.id + ")"));
                }
            }
            return Optional.empty();
        }
    }

    private static final class IdentificadoresUnicos implements RegraDeEntrada {
        @Override
        public Optional<ErroEntrada> verificar(List<Registro> registros) {
            List<String> vistos = new ArrayList<String>();
            for (Registro registro : registros) {
                if (vistos.contains(registro.id)) {
                    return Optional.of(new ErroEntrada(
                            "ID_DUPLICADO",
                            "identificador repetido na fila: " + registro.id));
                }
                vistos.add(registro.id);
            }
            return Optional.empty();
        }
    }

    private static final class TamanhoDosGrupos implements RegraDeEntrada {
        @Override
        public Optional<ErroEntrada> verificar(List<Registro> registros) {
            Map<String, Integer> contagem = new LinkedHashMap<String, Integer>();
            for (Registro registro : registros) {
                if (registro.grupo == null) {
                    continue;
                }
                Integer atual = contagem.get(registro.grupo);
                contagem.put(registro.grupo, atual == null ? 1 : atual + 1);
            }
            int maximo = Posicao.porEquipe();
            for (Map.Entry<String, Integer> entrada : contagem.entrySet()) {
                if (entrada.getValue().intValue() > maximo) {
                    return Optional.of(new ErroEntrada(
                            "GRUPO_MAIOR_QUE_EQUIPE",
                            "grupo " + entrada.getKey()
                                    + " possui " + entrada.getValue()
                                    + " jogadores, acima do tamanho maximo de uma equipe ("
                                    + maximo + ")"));
                }
            }
            return Optional.empty();
        }
    }

    private static final class Registro {
        final String id;
        final int habilidade;
        final String preferida;
        final List<String> alternativas;
        final int espera;
        final int latencia;
        final String grupo;
        final String regiao;

        Registro(String id, int habilidade, String preferida, List<String> alternativas,
                int espera, int latencia, String grupo, String regiao) {
            this.id = id;
            this.habilidade = habilidade;
            this.preferida = preferida;
            this.alternativas = alternativas;
            this.espera = espera;
            this.latencia = latencia;
            this.grupo = grupo;
            this.regiao = regiao;
        }
    }
}
