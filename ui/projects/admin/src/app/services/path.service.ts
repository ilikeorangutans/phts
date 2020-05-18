import { Injectable } from '@angular/core';

import { BasePathService } from 'projects/shared/src/app/services/base-path.service';

@Injectable({
  providedIn: 'root',
})
export class PathService {
  readonly api: string;

  constructor(private readonly basePath: BasePathService) {
    this.api = [this.basePath.apiHost, 'api', 'admin'].join('/');
  }

  joinToken(token: string): string {
    return [this.api, 'invite', token].join('/');
  }
}
