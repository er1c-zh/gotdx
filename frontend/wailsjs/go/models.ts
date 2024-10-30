export namespace api {
	
	export class StockMeta {
	    Market: number;
	    Code: string;
	    Desc: string;
	    Meta: proto.Security;
	
	    static createFrom(source: any = {}) {
	        return new StockMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Market = source["Market"];
	        this.Code = source["Code"];
	        this.Desc = source["Desc"];
	        this.Meta = this.convertValues(source["Meta"], proto.Security);
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
	    StockList: StockMeta[];
	    StockMap: {[key: string]: StockMeta};
	
	    static createFrom(source: any = {}) {
	        return new IndexInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IsConnected = source["IsConnected"];
	        this.Msg = source["Msg"];
	        this.StockCount = source["StockCount"];
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

}

export namespace proto {
	
	export class Security {
	    Code: string;
	    VolUnit: number;
	    Reserved1: number;
	    DecimalPoint: number;
	    Name: string;
	    PreClose: number;
	    Reserved2: number;
	
	    static createFrom(source: any = {}) {
	        return new Security(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Code = source["Code"];
	        this.VolUnit = source["VolUnit"];
	        this.Reserved1 = source["Reserved1"];
	        this.DecimalPoint = source["DecimalPoint"];
	        this.Name = source["Name"];
	        this.PreClose = source["PreClose"];
	        this.Reserved2 = source["Reserved2"];
	    }
	}

}

