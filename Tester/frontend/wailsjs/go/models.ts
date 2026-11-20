export namespace main {
	
	export class Response {
	    Success: boolean;
	    Data: any;
	    Error: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Success = source["Success"];
	        this.Data = source["Data"];
	        this.Error = source["Error"];
	    }
	}

}

