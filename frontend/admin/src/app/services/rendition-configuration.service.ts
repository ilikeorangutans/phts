import { Injectable } from '@angular/core';
import { Http } from '@angular/http';

import { PathService } from './path.service';

import { Collection } from '../models/collection';
import { RenditionConfiguration } from '../models/rendition-configuration';

@Injectable()
export class RenditionConfigurationService {

  constructor(
    private http: Http,
    private pathService: PathService
  ) { }

  forCollection(collection: Collection): Promise<RenditionConfiguration[]> {
    const p = this.pathService.renditionConfigurations(collection);
    console.log('RenditionConfigurationService::forCollection()', p);
    return this.http
      .get(p)
      .toPromise()
      .then((response) => {
        const configs = response.json() as PaginatedRenditionConfigurations;

        return configs.data;
      })
      .catch((e) => Promise.reject(e));
  }
}

interface PaginatedRenditionConfigurations {
  data: RenditionConfiguration[];
}
