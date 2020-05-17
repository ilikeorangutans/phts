import { Injectable } from '@angular/core';

import { BasePathService } from 'projects/shared/src/app/services/base-path.service';

@Injectable({
  providedIn: 'root',
})
export class PathService {
  readonly api: string;
  readonly share: string;

  constructor(private readonly basePath: BasePathService) {
    this.api = [this.basePath.apiHost, 'api'].join('/');
    this.share = [this.api, 'share'].join('/');
  }

  shareBySlug(slug: string): string {
    return [this.share, slug].join('/');
  }

  renditionBySlug(slug: string, renditionID: number): string {
    return [this.shareBySlug(slug), 'renditions', renditionID.toString()].join(
      '/'
    );
  }
}
