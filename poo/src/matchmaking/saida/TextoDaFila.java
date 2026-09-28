package matchmaking.saida;

import java.util.List;

public final class TextoDaFila {
    private TextoDaFila() {
    }

    public static String ids(List<String> ids) {
        if (ids.isEmpty()) {
            return "[]";
        }
        StringBuilder texto = new StringBuilder("[");
        for (int i = 0; i < ids.size(); i++) {
            if (i > 0) {
                texto.append(", ");
            }
            texto.append(ids.get(i));
        }
        texto.append("]");
        return texto.toString();
    }
}
