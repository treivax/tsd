CONTEXTE ET PRINCIPE :
Nous allons effectuer un changement majeur dans la structuration de tsd en séparant clairement les composantes que sont le moteur de règle RETE et le système de xuple-spaces (notre dérivation de tuple-spaces).
Pour l'instant les nœuds terminaux stockent les activations de règles (tokens) au lieu d'exécuter immédiatement les actions. Nous voulons rétablir un fonctionnement classique du moteur de règles où les actions sont déclanchées immédiatement. Nous allons définir et mettre en oeuvre, dans un premier temps, un ensemble minmal d'actions pour les règles. Certaines de ces règles permettront d'alimenter le système de xuple-spaces.

SPECIFICATION :
A) Au niveau RETE :
- Les actions sont déclanchées dès lors qu'un token déclancheur est produit.
- Un ensemble d'actions prédéfinies est utlisable par défaut. Il n'y a pas besoin de les déclarer (via la commande 'action') dans un programme tsd. Ces déclarations sont implicites et équivalentes ua parsing d'un programme tsd les déclarant (c'est même d'ailleurs une mise en oeuvre possible...). Les actions à définir et mettre en oeuvre sont :

    1. Print(string) qui affiche une chaîne de caractère sur la sortie standard,
    2. Log(string) qui génère une trace
    3. Update(fact) qui modifie un fait et met à jour les tokens qui lui sont liés dans le reseau RETE,
    4. Insert(fact) qui crée un nouveau fait et l'insère dans le réseau RETE,
    5. Retract(ID) qui supprime un fait (et tous les tokens qui luis sont liés) dans le réseau RETE,
    6. Xuple(xuple-name, fact)  qui crée un xuple dans le xuple-space dédié aux agents de nom 'name'.
Dans ces actions par défaut, un 'fact' est un type de fait quelconque parmi ceux qui auront été parsés par tsd au moment du déclanchement de la règle.


B) Au niveau xuple-spaces :
- Les xuples-spaces sont gérés par un nouveau module 'xuples' de tsd.
- Un xuple est un ensemble constitué : d'un fait et d'une liste de faits déclancheurs.
- Les xuples sont "consommés" par des agents externes.
- Un agent de nom xuple-name est le programme externe capable d'accéder aux xuples présents dans le xuple-space de nom xuple-name.
- Les xuples-spaces doivent être déclarés via la commande xuple-space de tsd. Cette commande définit le nom et la politique qui définit comment les xuples sont consomés par les agents.
- La politique d'un xuple-space exprime a minima :
   1. Comment est choisi un xuple particulier parmi l'esemble des xuples disposnibles dans le xuple-space lorsqu'un agent demande l'accès à un xuple : au hazard, en first-in/first-out ou en last-in/first-out.
   2. Comment sont consommés les xuples : le xuple ne peut être consommé qu'une fois, peut être consommé une fois par chaque agent lié au tuple space, un nombre limité de fois par chaque agent, etc.
   3. La durée de rétention d'un xuple : illimitée ou limitée à une durée avant que le xuple ne puisse plus être consommé.

Cette spécification (notamment celle des agents, du cycle de vie des xuples au sein des xuple-scpaces, les interactions entre les agants et tsd, etc.) sera affinée et complétée par la suite.

DEMANDE : 
Sans écrire de code, propose moi un plan d'action pour la mise en euvre :
- De la nouvelle commande xuple-space (au niveau du langage, du parsing, etc.)
- Des actions par défaut. Elles ne doivent pas être codées en dur dans le langage tsd. Par contre, la redéfinition d'une action définie par défaut doit bien provoquer une erreur de parsing exactement comme le fait de parser deux fois une même action. Une implémentation possible et simple serait de fait de parser réellement à l'initialisation un fichier de définition des actions par défaut.
- De la première version du nouveau module 'xuples'
- Des objets xuples. Un xuple est toujours créé au déclenchement d'une action 'Xuple', il est constitué du fait qui est passé en argument et des faits déclencheurs de l'action qui sont ceux présents dans le token combiné ayant déclanché l'action. Le xuple cré est stocké dans le xuples-space dont le nom est passé en argument.

OBLIGATOIRE : 
Crée ce plan d'action en respectant le prompt .github/common/develop.md. Le plan d'action doit obligatoirement avoir la forme d'une série ordonnée de prompts autonome et exécutables au sein d'une session unique de contexte de 128k. Chaque prompt de cette série est stocké dans un fichier qui lui est propre sein du sous-répertoire 'xuples' du répertoire 'scripts'. Tu ne peux en aucun cas synthétiser plusieurs prompts au sein d'un même fichier.
Les scripts du plan d'action doivent prendre obligatoirement en compte l'ensemble des règles et pratiques définies dans .github/prompts/common.md.
