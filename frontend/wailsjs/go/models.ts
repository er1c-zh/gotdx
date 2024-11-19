export namespace api {
	
	export enum MsgKey {
	    init = "init",
	    processMsg = "processMsg",
	    serverStatus = "serverStatus",
	}
	export class ExportStruct {
	    F0: models.ServerStatus;
	
	    static createFrom(source: any = {}) {
	        return new ExportStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.F0 = this.convertValues(source["F0"], models.ServerStatus);
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

}

