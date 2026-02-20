export namespace models {
	
	export class Photo {
	    id: number;
	    original_path: string;
	    library_path: string;
	    filename: string;
	    hash: string;
	    // Go type: time
	    date_taken: any;
	    camera_model: string;
	    latitude?: number;
	    longitude?: number;
	    // Go type: time
	    import_date: any;
	
	    static createFrom(source: any = {}) {
	        return new Photo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.original_path = source["original_path"];
	        this.library_path = source["library_path"];
	        this.filename = source["filename"];
	        this.hash = source["hash"];
	        this.date_taken = this.convertValues(source["date_taken"], null);
	        this.camera_model = source["camera_model"];
	        this.latitude = source["latitude"];
	        this.longitude = source["longitude"];
	        this.import_date = this.convertValues(source["import_date"], null);
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

