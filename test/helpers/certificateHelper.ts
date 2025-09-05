import { ethers } from 'hardhat';
import { BytesLike, HDNodeWallet } from 'ethers';
import { ChunkedX509CertStruct, BytesBlacklistStruct } from '../../typechain-types/contract/CertificateVerifier.sol/CertificateVerifier';

export class CertificateHelper {
    cert: ChunkedX509CertStruct;
    certHash: string;
    signer: HDNodeWallet;

    constructor(cert: ChunkedX509CertStruct, signer: HDNodeWallet, certHash: string) {
        this.cert = cert;
        this.signer = signer;
        this.certHash = certHash;
    }

    static genRandomWallet(): HDNodeWallet {
        return ethers.Wallet.createRandom();
    }

    static async gen(
        caCert: CertificateHelper | undefined = undefined,
        ca: string = '0x30030101ff',
        expirationDate?: number,
        userData?: string
    ): Promise<CertificateHelper> {
        const wallet = CertificateHelper.genRandomWallet();
        const publicKey = '0x' + wallet.signingKey.publicKey.slice(4);
        const serialNumber = ethers.keccak256(wallet.privateKey);
        const mrEnclave = ethers.keccak256(wallet.privateKey + '01');
        const mrSigner = ethers.keccak256(wallet.privateKey + '02');

        expirationDate = expirationDate ?? (await ethers.provider.getBlock("latest"))!.timestamp + 10000000000;
        userData = userData ?? CertificateHelper.genRandomWallet().privateKey;

        const cert: ChunkedX509CertStruct = {
            nonSerializedParts: ['0x06', '0x01', '0x02', '0x03', '0x04', '0x05', '0x07', '0x08', '0x09', '0x0a'],
            expirationDate: CertificateHelper.encodeDate(expirationDate),
            ca,
            userData,
            publicKey,
            serialNumber,
            mrEnclave,
            mrSigner,
            signature: '',
        };
        const certSigner = caCert === undefined ? wallet : caCert.signer;
        const { signature, certHash } = await CertificateHelper.signCert(cert, certSigner);
        cert.signature = signature;
        return new CertificateHelper(cert, wallet, certHash);
    }

    static encodeDate(unixTimestamp: number): string {
        const date = new Date(unixTimestamp * 1000).toISOString();
        const regExp = /-|T|:/g;
        return ethers.hexlify(ethers.toUtf8Bytes(date.replace(regExp, '')));
    }

    static decodeDate(dateStr: BytesLike): number {
        dateStr = ethers.toUtf8String(dateStr);
        let year;
        let startIndex;
        if (dateStr.length === 13) {
            startIndex = 2; // UTCTime format (YYMMDDhhmmssZ)
            year = parseInt(dateStr.slice(0, startIndex), 10);
            year += year < 50 ? 2000 : 1900;
        } else if (dateStr.length === 19) {
            startIndex = 4; // GeneralizedTime format (YYYYMMDDhhmmss.mmmZ)
            year = parseInt(dateStr.slice(0, 4), 10);
        } else {
            throw new Error(`Invalid date format: ${dateStr}`);
        }
        const month = parseInt(dateStr.slice(startIndex, startIndex + 2), 10) - 1; // Months are 0-indexed
        const day = parseInt(dateStr.slice(startIndex + 2, startIndex + 4), 10);
        const hours = parseInt(dateStr.slice(startIndex + 4, startIndex + 6), 10);
        const minutes = parseInt(dateStr.slice(startIndex + 6, startIndex + 8), 10);
        const seconds = parseInt(dateStr.slice(startIndex + 8, startIndex + 10), 10);
        const date = new Date(Date.UTC(year, month, day, hours, minutes, seconds));
        const unixTimestamp = Math.floor(date.getTime() / 1000);
        return unixTimestamp;
    }

    static async signCert(cert: ChunkedX509CertStruct, user: HDNodeWallet): Promise<{ signature: string; certHash: string }> {
        const hash = ethers.solidityPackedSha256(
            ['bytes', 'bytes', 'bytes', 'bytes', 'bytes', 'bytes', 'bytes', 'bytes5', 'bytes', 'bytes32', 'bytes', 'bytes32', 'bytes', 'bytes32'],
            [
                cert.nonSerializedParts[1],
                cert.serialNumber,
                cert.nonSerializedParts[2],
                cert.expirationDate,
                cert.nonSerializedParts[3],
                cert.publicKey,
                cert.nonSerializedParts[4],
                cert.ca,
                cert.nonSerializedParts[5],
                cert.userData,
                cert.nonSerializedParts[6],
                cert.mrEnclave,
                cert.nonSerializedParts[7],
                cert.mrSigner,
            ]
        );

        const signature = await user.signingKey.sign(hash);
        const rawSignature = signature.r + signature.s.slice(2);

        return { signature: rawSignature, certHash: hash };
    }
}

export class BlacklistHelper {
    static gen(): BytesBlacklistStruct[] {
        const result: BytesBlacklistStruct[] = [];
        const count = getRandomInt(1, 10);
        for (let i = 0; i < count; i++) {
            result.push({
                item: CertificateHelper.genRandomWallet().privateKey,
                isBlacklisted: getRandomInt(0, 1) === 1,
            });
        }
        return result;
    }
}

function getRandomInt(min: number, max: number) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}