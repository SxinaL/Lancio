export namespace database {
	
	export class Word {
	    id: number;
	    word: string;
	    phonetic: string;
	    translation: string;
	    part_of_speech: string;
	    example_sentence: string;
	    difficulty: number;
	
	    static createFrom(source: any = {}) {
	        return new Word(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.word = source["word"];
	        this.phonetic = source["phonetic"];
	        this.translation = source["translation"];
	        this.part_of_speech = source["part_of_speech"];
	        this.example_sentence = source["example_sentence"];
	        this.difficulty = source["difficulty"];
	    }
	}

}

export namespace services {
	
	export class WordDetail {
	    word: string;
	    phonetic_us: string;
	    phonetic_uk: string;
	    translation: string;
	    part_of_speech: string;
	    definitions: string[];
	    examples: string[];
	    synonyms: string[];
	    antonyms: string[];
	
	    static createFrom(source: any = {}) {
	        return new WordDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.word = source["word"];
	        this.phonetic_us = source["phonetic_us"];
	        this.phonetic_uk = source["phonetic_uk"];
	        this.translation = source["translation"];
	        this.part_of_speech = source["part_of_speech"];
	        this.definitions = source["definitions"];
	        this.examples = source["examples"];
	        this.synonyms = source["synonyms"];
	        this.antonyms = source["antonyms"];
	    }
	}

}

