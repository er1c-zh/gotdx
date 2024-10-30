export namespace main {
	
	export class Status {
	    Msg: string;
	
	    static createFrom(source: any = {}) {
	        return new Status(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Msg = source["Msg"];
	    }
	}

}

