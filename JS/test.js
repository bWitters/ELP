const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
}); // Définie l'objet rl qui lit les entrées du terminale et envoye les sorties dans le terminal

const wordToGuess = "piano"; // Mot mystère à deviner
let hints = [];

console.log("Bienvenue dans Just One - Version Console");
console.log("Un joueur doit deviner un mot, les autres donnent un indice.");
console.log("Si plusieurs joueurs donnent le même indice, il est supprimé !\n");

function askHint() {
    rl.question("Entrez un indice (ou tapez 'fin' pour arrêter) : ", (hint) => {
        if (hint.toLowerCase() === 'fin') {
            processHints();
        } else {
            hints.push(hint.toLowerCase());
            askHint();
        }
    }); // Fonction qui attend d'avoir une entrée pour s'executer
}

function processHints() {
    const hintCount = {};
    hints.forEach(hint => {
        hintCount[hint] = (hintCount[hint] || 0) + 1; // Pas d'erreur si la clef n'existe pas 0_0, Cette ligne compte les occurences d'un mot
    });

    const uniqueHints = Object.keys(hintCount).filter(hint => hintCount[hint] === 1);
    // Object.keys() c'est les clefs du dictionnaire hintCount.
    // Array.filter() ici filtre sur le fait que pour un hint, la valeur dans le dictionnaire soit égale à 1.

    
    console.log("\nIndices valides :", uniqueHints.length > 0 ? uniqueHints.join(", ") : "Aucun (tous supprimés)");
    // Test, si la taille de la liste uniqueHints > 0 alors on affiche les elements de la liste séparé par des virgules
    // Sinon on affiche Aucun (tous supprimés).

    rl.question("Devinez le mot mystère : ", (answer) => {
        if (answer.toLowerCase() === wordToGuess) {
            console.log("Bravo, vous avez trouvé le mot !");
        } else {
            console.log("Dommage, le mot était :", wordToGuess);
        }
        rl.close(); // Ferme l'objet qui lit et ecrit sur le terminal.
    }); // Pareil que l'autre quesion, attend une entrée dans le terminale pour proceder.
}

askHint();
