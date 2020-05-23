import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';

import { Observable } from 'rxjs';
import { first, map } from 'rxjs/operators';

import { PathService } from './path.service';
import { Collection } from '../models/collection';
import { RenditionConfiguration } from '../models/rendition-configuration';

@Injectable({ providedIn: 'root' })
export class RenditionConfigurationService {
  constructor(private http: HttpClient, private pathService: PathService) {}

  forCollection(collection: Collection): Observable<RenditionConfiguration[]> {
    const p = this.pathService.renditionConfigurations(collection);
    return this.http.get<PaginatedRenditionConfigurations>(p).pipe(
      map((response) => {
        return response.data;
      }),
      first()
    );
  }

  save(
    collection: Collection,
    config: RenditionConfiguration
  ): Promise<RenditionConfiguration> {
    const p = this.pathService.renditionConfigurations(collection);
    return this.http.post<RenditionConfiguration>(p, config).toPromise();
  }
}

interface PaginatedRenditionConfigurations {
  data: RenditionConfiguration[];
}
