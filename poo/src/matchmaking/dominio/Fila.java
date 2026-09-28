package matchmaking.dominio;

import java.util.ArrayList;
import java.util.Collection;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

/** Agregado da fila: dono dos jogadores e dos grupos, na ordem de chegada. */
public final class Fila {
    private final List<Jogador> jogadores;
    private final List<Grupo> grupos;

    private Fila(List<Jogador> jogadores, List<Grupo> grupos) {
        this.jogadores = jogadores;
        this.grupos = grupos;
    }

    public static Fila formar(List<Inscricao> inscricoes) {
        Map<String, Grupo> porId = new LinkedHashMap<String, Grupo>();
        List<Jogador> criados = new ArrayList<Jogador>();
        for (Inscricao inscricao : inscricoes) {
            Grupo grupo = null;
            if (inscricao.grupoId() != null) {
                grupo = porId.get(inscricao.grupoId());
                if (grupo == null) {
                    grupo = new Grupo(inscricao.grupoId());
                    porId.put(inscricao.grupoId(), grupo);
                }
            }
            Jogador jogador = Jogador.aPartirDe(inscricao, grupo);
            if (grupo != null) {
                grupo.adicionar(jogador);
            }
            criados.add(jogador);
        }
        for (Grupo grupo : porId.values()) {
            grupo.fechar();
        }
        return new Fila(List.copyOf(criados), List.copyOf(porId.values()));
    }

    public List<Jogador> jogadores() {
        return jogadores;
    }

    public List<Grupo> grupos() {
        return grupos;
    }

    public int tamanho() {
        return jogadores.size();
    }

    public boolean contem(Jogador jogador) {
        for (Jogador atual : jogadores) {
            if (atual.equals(jogador)) {
                return true;
            }
        }
        return false;
    }

    public List<String> ids() {
        List<String> ids = new ArrayList<String>();
        for (Jogador jogador : jogadores) {
            ids.add(jogador.id());
        }
        return ids;
    }

    public List<String> idsExceto(Collection<Jogador> selecionados) {
        Set<String> escolhidos = new HashSet<String>();
        for (Jogador jogador : selecionados) {
            escolhidos.add(jogador.id());
        }
        List<String> restante = new ArrayList<String>();
        for (Jogador jogador : jogadores) {
            if (!escolhidos.contains(jogador.id())) {
                restante.add(jogador.id());
            }
        }
        return restante;
    }
}
