export namespace entity {
	
	export class Vocabulary {
	    id: number;
	    word: string;
	    phonetic: string;
	    translation: string;
	    example_sentence: string;
	    vocabulary_bank_id: number;
	
	    static createFrom(source: any = {}) {
	        return new Vocabulary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.word = source["word"];
	        this.phonetic = source["phonetic"];
	        this.translation = source["translation"];
	        this.example_sentence = source["example_sentence"];
	        this.vocabulary_bank_id = source["vocabulary_bank_id"];
	    }
	}
	export class VocabularyBank {
	    id: number;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new VocabularyBank(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}

}

export namespace handler {
	
	export class MinSize {
	    WindowSizeMinW: number;
	    WindowSizeMinH: number;
	
	    static createFrom(source: any = {}) {
	        return new MinSize(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WindowSizeMinW = source["WindowSizeMinW"];
	        this.WindowSizeMinH = source["WindowSizeMinH"];
	    }
	}

}

