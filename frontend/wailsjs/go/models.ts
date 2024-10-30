export namespace main {
	
	export class IndexInfo {
	    IsConnected: boolean;
	    Msg: string;
	    StockCount: number;
	
	    static createFrom(source: any = {}) {
	        return new IndexInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsConnected = source["IsConnected"];
	        this.Msg = source["Msg"];
	        this.StockCount = source["StockCount"];
	    }
	}

}

