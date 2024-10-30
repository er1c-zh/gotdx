export namespace main {
	
	export class StockMeta {
	    Market: number;
	    Code: string;
	    Desc: string;
	
	    static createFrom(source: any = {}) {
	        return new StockMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Market = source["Market"];
	        this.Code = source["Code"];
	        this.Desc = source["Desc"];
	    }
	}
	export class StockMarketList {
	    Market: number;
	    MarketStr: string;
	    Count: number;
	    StockList: StockMeta[];
	    StockMap: {[key: string]: StockMeta};
	
	    static createFrom(source: any = {}) {
	        return new StockMarketList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Market = source["Market"];
	        this.MarketStr = source["MarketStr"];
	        this.Count = source["Count"];
	        this.StockList = this.convertValues(source["StockList"], StockMeta);
	        this.StockMap = this.convertValues(source["StockMap"], StockMeta, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class IndexInfo {
	    IsConnected: boolean;
	    Msg: string;
	    StockCount: number;
	    AllStock: StockMarketList[];
	
	    static createFrom(source: any = {}) {
	        return new IndexInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsConnected = source["IsConnected"];
	        this.Msg = source["Msg"];
	        this.StockCount = source["StockCount"];
	        this.AllStock = this.convertValues(source["AllStock"], StockMarketList);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

