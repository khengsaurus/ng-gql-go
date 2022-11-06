import { Injectable } from '@angular/core';
import { UserService } from './user.service';
import { environment } from 'src/environments/environment';
import { HttpClient } from '@angular/common/http';
import { of, switchMap, tap } from 'rxjs';

const config = { headers: { Authorization: 'Bearer -' } };
const MAX_MB = 10;

@Injectable({ providedIn: 'root' })
export class FilesService {
  constructor(private http: HttpClient, private userService: UserService) {}
  private api_route = `${environment.restApi}/files`;

  getSignedPutURL$(todoId: string, file: File) {
    const userId = this.userService.currentUser?.id;
    if (!todoId || !userId || !file.name) return null;
    if (file.size > MAX_MB * 1000_000) {
      console.error(`File size cannot exceed ${MAX_MB}MB`);
      return null;
    }
    return this.http.post(
      this.api_route,
      { userId, todoId, fileName: file.name },
      config
    );
  }

  /**
   * Upload a file to S3 via presigned PUT url, and return the file key
   */
  uploadFile$(todoId: string, file: File) {
    const presignReq = this.getSignedPutURL$(todoId, file);
    return presignReq
      ? presignReq.pipe(
          switchMap((presignRes: any) => {
            if (presignRes?.key && presignRes?.url) {
              return fetch(presignRes.url, {
                method: 'put',
                headers: environment.production
                  ? { 'Content-Type': 'multipart/form-data' }
                  : {
                      'Content-Type': 'application/json',
                      'x-amz-acl': 'public-read-write',
                    },
                body: file,
              })
                .then((uploadRes) => {
                  if (uploadRes?.status === 200) {
                    return presignRes.key;
                  } else {
                    throw new Error(
                      `PUT request to presigned S3 URL returned status ${uploadRes?.status}`
                    );
                  }
                })
                .catch((err) => {
                  console.error(err);
                  return '';
                });
            } else {
              return '';
            }
          })
        )
      : of('');
  }
}
