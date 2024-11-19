export namespace api {
	
	export enum MsgKey {
	    init = "init",
	    processMsg = "processMsg",
	    serverStatus = "serverStatus",
	    subscribeBroadcast = "subscribeBroadcast",
	}
	export class ExportStruct {
	    F0: models.ServerStatus;
	    F1: v2.RealtimeRespItem[];
	
	    static createFrom(source: any = {}) {
	        return new ExportStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.F0 = this.convertValues(source["F0"], models.ServerStatus);
	        this.F1 = this.convertValues(source["F1"], v2.RealtimeRespItem);
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
	export class SubscribeReq {
	    Group: string;
	    Code: string[];
	    QuoteType: string;
	
	    static createFrom(source: any = {}) {
	        return new SubscribeReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Group = source["Group"];
	        this.Code = source["Code"];
	        this.QuoteType = source["QuoteType"];
	    }
	}

}

export namespace models {
	
	export enum MarketType {
	    深圳 = 0x0,
	    上海 = 0x1,
	    北京 = 0x2,
	}
	export class ProcessInfo {
	    Type: number;
	    Msg: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Type = source["Type"];
	        this.Msg = source["Msg"];
	    }
	}
	export class ServerStatus {
	    Connected: boolean;
	    ServerInfo: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Connected = source["Connected"];
	        this.ServerInfo = source["ServerInfo"];
	    }
	}
	export class StockMetaItem {
	    Market: MarketType;
	    Code: string;
	    Desc: string;
	    PinYinInitial: string;
	
	    static createFrom(source: any = {}) {
	        return new StockMetaItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Market = source["Market"];
	        this.Code = source["Code"];
	        this.Desc = source["Desc"];
	        this.PinYinInitial = source["PinYinInitial"];
	    }
	}

}

export namespace v2 {
	
	export class RealtimeRespItem {
	    Market: number;
	    Code: string;
	    CurrentPrice: number;
	    YesterdayCloseDelta: number;
	    OpenDelta: number;
	    HighDelta: number;
	    LowDelta: number;
	    TotalVolume: number;
	    CurrentVolume: number;
	    TotalAmount: number;
	
	    static createFrom(source: any = {}) {
	        return new RealtimeRespItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Market = source["Market"];
	        this.Code = source["Code"];
	        this.CurrentPrice = source["CurrentPrice"];
	        this.YesterdayCloseDelta = source["YesterdayCloseDelta"];
	        this.OpenDelta = source["OpenDelta"];
	        this.HighDelta = source["HighDelta"];
	        this.LowDelta = source["LowDelta"];
	        this.TotalVolume = source["TotalVolume"];
	        this.CurrentVolume = source["CurrentVolume"];
	        this.TotalAmount = source["TotalAmount"];
	    }
	}

}

