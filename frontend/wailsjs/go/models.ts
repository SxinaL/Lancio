export namespace entity {
	
	export class Vocabulary {
	    id: number;
	    word: string;
	    phonetic: string;
	    translation: string;
	    example_sentence: string;
	    vocabulary_bank_id: number;
	    review_count: number;
	    review_status: string;
	
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
	        this.review_count = source["review_count"];
	        this.review_status = source["review_status"];
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
	
	export class APIConfigKey {
	    appid: string;
	    key: string;
	
	    static createFrom(source: any = {}) {
	        return new APIConfigKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appid = source["appid"];
	        this.key = source["key"];
	    }
	}
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

